package files

import "io"

type Progress func(int)

func ProgressWriter(w io.Writer, report func(int)) io.Writer {
	return io.MultiWriter(w, Progress(report))
}

func ProgressReader(r io.Reader, report func(int)) io.Reader {
	return io.MultiReader(r, Progress(report))
}

func (p Progress) Write(buf []byte) (n int, err error) { return p.Read(buf) }
func (p Progress) Read(buf []byte) (n int, err error)  { n = len(buf); p(n); return }
