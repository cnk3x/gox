package strs

import (
	"strconv"
	"time"
)

var (
	dunit_h = [3]string{"d", "h", "m"}
	dbase_h = [3]time.Duration{time.Hour * 24, time.Hour, time.Minute}
	dunit_l = [4]string{"s", "ms", "µs", "ns"}
	dbase_l = [4]time.Duration{time.Second, time.Millisecond, time.Microsecond, 1}
)

func Duration(t time.Duration) string {
	if t == 0 {
		return ""
	}

	r := make([]string, 0, 7)
	if t < 0 {
		r = append(r, "-")
		t = -t
	}

	if t >= time.Minute { // "d", "h", "m"
		for i := 0; i < len(dbase_h) && t > 0; i++ {
			if base := dbase_h[i]; t >= base {
				r = append(r, strconv.Itoa(int(t/base)), dunit_h[i])
				t = t % base
			}
		}
	} else { // "s", "ms", "µs", "ns"
		for i := 0; i < len(dbase_l) && t > 0; i++ {
			if base := dbase_l[i]; t >= base {
				r = append(r, strconv.FormatFloat(float64(t)/float64(base), 'f', 3, 32)+dunit_l[i])
				break
			}
		}
	}

	return Join(r, "")
}

func ParseDuration(s string) (t time.Duration, err error) {
	if s == "" {
		return
	}

	var d string
	var f bool

	if f := s[0] == '-'; f {
		s = s[1:]
	}

	i := Index(s, "d")
	if i > 0 {
		d, s = s[:i], s[i+1:]
	}

	if t, err = time.ParseDuration(s); err != nil {
		return
	}

	if i > 0 {
		n, e := strconv.Atoi(d)
		if err = e; e != nil {
			return
		}
		t += time.Duration(n) * time.Hour * 24
	}

	if f {
		t = -t
	}

	return
}
