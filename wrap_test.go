package structlog_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/structlog"
)

func TestWrapErr(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf)
	err := log.WrapErr(io.EOF, "a", 10, "b", 20)
	err = log.WrapErr(err, "a", 11, "c", 30)
	err = log.WrapErr(err, "a", 12, "d", 40)
	log.Warn("hmm", "c", 31, "e", 50, "err", err)
	t.Match(buf.String(), "`hmm` a=12 b=20 c=31 d=40 e=50 err=EOF")
}
