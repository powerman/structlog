package structlog_test

import (
	"bytes"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/structlog"
)

func testPanic(log *structlog.Logger) {
	defer log.Recover(nil)
	panic("oops")
}

//go:noinline
func testPanicWrapper(log *structlog.Logger) {
	testPanic(log)
	log.Info("wrapper")
}

//go:inline
func testPanicThinWrapper(log *structlog.Logger) {
	testPanic(log)
}

func TestRecover(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf)
	testPanicAnon := func(log *structlog.Logger) {
		defer log.Recover(nil)
		panic("oops")
	}
	testPanicAnon(log)
	t.Match(buf.String(), `@ structlog_test.TestRecover.func1\(caller_test.go:34\)`)
	buf.Reset()
	testPanic(log)
	t.Match(buf.String(), `@ structlog_test.testPanic\(caller_test.go:14\)`)
	buf.Reset()
	testPanicWrapper(log)
	t.Match(buf.String(), `@ structlog_test.testPanic\(caller_test.go:14\)`)
	buf.Reset()
	testPanicThinWrapper(log)
	t.Match(buf.String(), `@ structlog_test.testPanic\(caller_test.go:14\)`)
}
