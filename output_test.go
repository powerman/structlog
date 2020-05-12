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

func (bp *bufPrinter) Print(v ...interface{}) { fmt.Fprint(&bp.Buffer, append(v, "\n")...) }

func TestDefaultPrinter(tt *testing.T) {
	t := check.T(tt)
	defer stdlog.SetOutput(os.Stderr)
	var buf bytes.Buffer
	stdlog.SetOutput(&buf)
	log := structlog.New(structlog.KeyUnit, "structlog")
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestDefaultPrinter(output_test.go:24)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestDefaultPrinter(output_test.go:25)\n")
}

func TestPrinter(tt *testing.T) {
	t := check.T(tt)
	var buf bufPrinter
	log := structlog.New(structlog.KeyUnit, "structlog").SetPrinter(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestPrinter(output_test.go:35)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestPrinter(output_test.go:36)\n")
}

func TestOutput(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	log := structlog.New(structlog.KeyUnit, "structlog").SetOutput(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf "+unit+": `something happens` k1=v1 k2=v2 \t@ structlog_test.TestOutput(output_test.go:46)\n"+
		"structlog.test["+pid+"] WRN "+unit+": `oops` \t@ structlog_test.TestOutput(output_test.go:47)\n")
}
