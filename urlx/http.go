package urlx

import (
	"cmp"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/cnk3x/gox/files"
	"github.com/cnk3x/gox/strs"
)

type (
	RequestOptions struct {
		method    string
		url       string
		body      func(context.Context) (io.ReadCloser, string, error)
		responses []func(*http.Response) error
		headers   []func(header http.Header)
		clients   []func(*http.Client)
	}
	RequestOption   func(*RequestOptions)
	ResponseProcess func(*http.Response) error
)

func Request(ctx context.Context, options ...RequestOption) (err error) {
	rOpts := RequestOptions{method: http.MethodGet}

	for _, opt := range options {
		opt(&rOpts)
	}

	var (
		req       *http.Request
		resp      *http.Response
		input     io.ReadCloser
		inputMime string
	)

	if rOpts.body != nil {
		if input, inputMime, err = rOpts.body(ctx); err != nil {
			return
		}
		if input != nil {
			defer input.Close()
		}
		rOpts.method = cmp.Or(rOpts.method, http.MethodPost)
	} else {
		rOpts.method = cmp.Or(rOpts.method, http.MethodGet)
	}

	if req, err = http.NewRequestWithContext(ctx, rOpts.method, rOpts.url, input); err != nil {
		return
	}

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	for _, headerFn := range rOpts.headers {
		headerFn(req.Header)
	}

	if inputMime != "" {
		req.Header.Set("Content-Type", inputMime)
	}

	client := &http.Client{Transport: DefaultTransport}
	for _, clientFn := range rOpts.clients {
		clientFn(client)
	}

	if resp, err = client.Do(req); err != nil {
		return
	}
	output := resp.Body
	defer output.Close()

	for _, responseFn := range rOpts.responses {
		if err = responseFn(resp); err != nil {
			return
		}
	}

	return
}

func Url(url string, method ...string) RequestOption {
	return func(ro *RequestOptions) {
		ro.url = url
		if len(method) > 0 {
			ro.method = method[0]
		}
	}
}

func Input(input func(context.Context) (body io.ReadCloser, contentType string, err error)) RequestOption {
	return func(ro *RequestOptions) {
		ro.body = input
	}
}

func Header(headerFunc func(header http.Header)) RequestOption {
	return func(ro *RequestOptions) {
		ro.headers = append(ro.headers, headerFunc)
	}
}

func HeaderSet(n, v string) RequestOption {
	return Header(func(header http.Header) { header[n] = []string{v} })
}

func HeaderSets(nv ...string) RequestOption {
	return Header(func(header http.Header) {
		for i := 0; i < len(nv)-1; i += 2 {
			header[nv[i]] = []string{nv[i+1]}
		}
	})
}

func HeaderAdd(n, v string) RequestOption {
	return Header(func(header http.Header) { header[n] = append(header[n], v) })
}

func HeaderAdds(nv ...string) RequestOption {
	return Header(func(header http.Header) {
		for i := 0; i < len(nv)-1; i += 2 {
			header[nv[i]] = append(header[nv[i]], nv[i+1])
		}
	})
}

func HeaderDel(ns ...string) RequestOption {
	return Header(func(header http.Header) {
		for _, n := range ns {
			delete(header, n)
		}
	})
}

func Headers(lines ...string) RequestOption {
	return Header(func(header http.Header) {
		for _, line := range lines {
			if k, v, ok := strs.Cut(line, ":"); ok {
				k, v = strs.TrimSpace(k), strs.TrimSpace(v)
				header[k] = append(header[k], v)
			}
		}
	})
}

func Process(response ResponseProcess) RequestOption {
	return func(ro *RequestOptions) {
		ro.responses = append(ro.responses, response)
	}
}

func SaveResponse(to string, reportFn ...func(cur, total int64)) RequestOption {
	return Process(func(resp *http.Response) (err error) {
		temp := to + ".downloading_temp"
		defer os.Remove(temp)

		current, total := int64(0), resp.ContentLength

		var reports []func(n int)
		if len(reportFn) > 0 {
			reports = append(reports, func(n int) { atomic.AddInt64(&current, int64(n)) })
			for _, report := range reportFn {
				report := report
				reports = append(reports, func(n int) { report(current, total) })
			}
		}

		if err = files.SaveFile(temp, resp.Body, reports...); err != nil {
			return
		}

		if err = os.Remove(to); os.IsNotExist(err) || err == nil {
			err = os.Rename(temp, to)
		}
		return
	})
}

func GetBytes(data *[]byte) RequestOption {
	return Process(func(resp *http.Response) (err error) {
		defer resp.Body.Close()
		*data, err = io.ReadAll(resp.Body)
		return
	})
}

func ReadBytes(read func([]byte) error) RequestOption {
	return Process(func(resp *http.Response) (err error) {
		defer resp.Body.Close()
		var data []byte
		if data, err = io.ReadAll(resp.Body); err == nil {
			err = read(data)
		}
		return
	})
}

// func ReadJSON(read func(gjson.Result) error) RequestOption {
// 	return ReadBytes(func(b []byte) (err error) {
// 		return read(gjson.ParseBytes(b))
// 	})
// }

func Client(clientFn func(*http.Client)) RequestOption {
	return func(ro *RequestOptions) {
		ro.clients = append(ro.clients, clientFn)
	}
}

var DefaultTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	DialContext:           (&net.Dialer{Timeout: 15 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}
