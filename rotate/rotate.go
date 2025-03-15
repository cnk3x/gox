package rotate

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Options struct {
	Path       string `json:"path,omitempty" yaml:"path,omitempty"`               // 文件路径
	MaxSize    int64  `json:"max_size,omitempty" yaml:"max_size,omitempty"`       // 单文件最大大小
	MaxBackups int    `json:"max_backups,omitempty" yaml:"max_backups,omitempty"` // 最大备份文件数量
}

func New(options Options) io.WriteCloser {
	w := &Writer{
		Path:       options.Path,
		MaxSize:    options.MaxSize,
		MaxBackups: options.MaxBackups,
	}
	w.Init()
	return w
}

type Writer struct {
	Path       string // 文件路径
	MaxSize    int64  // 单文件最大大小
	MaxBackups int    // 最大备份文件数量

	cur  *os.File
	size atomic.Int64

	dir  string
	name string
	ext  string

	tt *time.Timer
	tr atomic.Bool

	mu sync.Mutex
}

func (w *Writer) Init() {
	w.dir, w.name = filepath.Split(w.Path)
	w.ext = filepath.Ext(w.name)
	w.name = w.name[:len(w.name)-len(w.ext)]
	// w.sync = TimerFunc(func() { _ = w.cur.Sync() })

	w.tt = time.AfterFunc(time.Second, func() {
		w.cur.Sync()
		w.tr.Store(false)
	})

	_ = os.MkdirAll(filepath.Dir(w.Path), 0o777)
	w.create()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if n, err = w.cur.Write(p); err != nil {
		return
	}

	if x := w.size.Add(int64(n)); x >= w.MaxSize {
		w.mu.Lock()
		err = w.rotate()
		w.mu.Unlock()
		if err != nil {
			return
		}
	} else if w.tr.CompareAndSwap(false, true) {
		w.tt.Reset(time.Second)
	}
	return
}

func (w *Writer) Close() error {
	w.tt.Stop()
	w.tr.Store(false)
	return w.cur.Close()
}

func (w *Writer) create() (err error) {
	if w.cur, err = os.OpenFile(w.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666); err != nil {
		return
	}
	if stat, _ := w.cur.Stat(); stat != nil {
		w.size.Store(stat.Size())
	} else {
		w.size.Store(0)
	}
	return
}

func (w *Writer) backup() {
	if err := w.gzBackup(w.Path+".backup", w.nameNow()+".gz", true); err == nil {
		if w.MaxBackups > 0 {
			files, _ := filepath.Glob(filepath.Join(w.dir, w.name+"*"+w.ext+".gz"))
			if l := len(files); l > w.MaxBackups {
				sort.Strings(files)
				for i := w.MaxBackups; i < l; i++ {
					os.Remove(files[i])
				}
			}
		}
	}
	return
}

func (w *Writer) rotate() (err error) {
	if w.cur != nil {
		_ = w.cur.Sync()

		if err = closeIt(w.cur, nil); err != nil {
			return
		}

		backupPath := w.Path + ".backup"
		if err = os.Rename(w.Path, backupPath); err != nil {
			return
		}

		go w.backup()
	}

	return w.create()
}

func (w *Writer) nameNow() string {
	return filepath.Join(w.dir, w.name+"-"+time.Now().Format("20060102-150405")+w.ext)
}

func (w *Writer) gzBackup(sourcePath string, targetPath string, removeSource bool) (err error) {
	var src, dst *os.File

	if src, err = os.Open(sourcePath); err != nil {
		return
	}

	err = func() (err error) {
		if dst, err = os.Create(targetPath); err != nil {
			return
		}
		gw := gzip.NewWriter(dst)
		_, err = io.Copy(gw, src)
		err = closeIt(gw, err)
		err = closeIt(dst, err)
		if err != nil {
			os.Remove(targetPath)
		}
		return
	}()

	if err = closeIt(src, err); err != nil {
		return
	}

	if removeSource {
		if err = os.Remove(sourcePath); err != nil {
			return
		}
	}

	return
}

func closeIt(c io.Closer, err error) error {
	if c != nil {
		if e := c.Close(); err == nil && e != nil {
			err = e
		}
	}
	return err
}
