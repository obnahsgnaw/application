package writer

import (
	"io"
)

type MultiWriter struct {
	w []io.Writer
}

func NewMultiWriter(w ...io.Writer) io.Writer {
	if len(w) == 0 {
		return nil
	}
	var ww []io.Writer
	for _, o := range w {
		if o != nil {
			ww = append(ww, o)
		}
	}
	return &MultiWriter{
		w: ww,
	}
}

func (s *MultiWriter) Write(p []byte) (int, error) {
	for _, w := range s.w {
		if n, err := w.Write(p); err != nil {
			return n, err
		}
	}
	return len(p), nil
}
