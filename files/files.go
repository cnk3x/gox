package files

import (
	"io"
	"os"
)

func SaveFile(dstFile string, src io.Reader, onProgress ...func(n int)) (err error) {
	dst, e := os.Create(dstFile)
	if err = e; err == nil {
		defer dst.Close()
		var w io.Writer
		if len(onProgress) > 0 {
			w = io.MultiWriter(dst, Progress(func(n int) {
				for _, p := range onProgress {
					if p != nil {
						p(n)
					}
				}
			}))
		} else {
			w = dst
		}
		_, err = io.Copy(w, src)
	}
	return
}
