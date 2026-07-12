package structlog_test

import (
	"bytes"
	"fmt"
	stdlog "log"
	"os"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/structlog"
)

type bufPrinter struct{ bytes.Buffer }

func (bp *bufPrinter) Print(v ...any) { fmt.Fprint(&bp.Buffer, append(v, "\n")...) }

func TestDefaultPrinter(tt *testing.T) {
	tt.Parallel()
	t := check.Must(tt)
	t.Cleanup(func() { stdlog.SetOutput(os.Stderr) })
	var buf bytes.Buffer
	stdlog.SetOutput(&buf)
	log := structlog.New()
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestDefaultPrinter(output_test.go:26)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestDefaultPrinter(output_test.go:27)\n")
}

func TestPrinter(tt *testing.T) {
	tt.Parallel()
	t := check.Must(tt)
	var buf bufPrinter
	log := structlog.New().SetPrinter(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestPrinter(output_test.go:38)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestPrinter(output_test.go:39)\n")
}

func TestOutput(tt *testing.T) {
	tt.Parallel()
	t := check.Must(tt)
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestOutput(output_test.go:50)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestOutput(output_test.go:51)\n")
}
