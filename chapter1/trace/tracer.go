package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

// io.Writerで受け付けるということは、ユーザーが出力先を自由に選べるということを意味する
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
