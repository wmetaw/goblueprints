package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {

	var buf bytes.Buffer

	tracer := New(&buf)

	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは trace package")
		if buf.String() != "こんにちは trace package\n" {
			t.Errorf("'%s'という誤った文字列が検出されました", buf.String())
		}
	}
}
