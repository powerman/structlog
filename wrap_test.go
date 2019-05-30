package structlog_test

import (
	"bytes"
	"io"
	"os"
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

	t.Nil(log.WrapErr(nil, "a", 10))
}

func ExampleLogger_WrapErr() {
	// Use NewZeroLogger to avoid reconfiguring
	// structlog.DefaultLogger in example, but in real code usually
	// reconfiguring DefaultLogger is better than using NewZeroLogger.
	log := structlog.NewZeroLogger().
		SetOutput(os.Stdout).
		SetPrefixKeys(structlog.KeyLevel).
		SetKeysFormat(map[string]string{
			structlog.KeyLevel:   "%[2]s",
			structlog.KeyMessage: " %#[2]q",
		})

	lowLevelFunc := func() error {
		return log.WrapErr(io.EOF, "details", "about error")
	}
	middleLevelFunc := func(action string) error {
		if err := lowLevelFunc(); err != nil {
			return log.WrapErr(err, "action", action)
		}
		return nil
	}
	topLevelFunc := func() {
		if err := middleLevelFunc("doit"); err != nil {
			log.Warn("log only at top level", "err", err)
		}
	}
	topLevelFunc()
	// Output:
	// WRN `log only at top level` details=about error action=doit err=EOF
}
