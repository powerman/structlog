package structlog_test

import (
	"bytes"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/structlog"
)

type bufPrinter struct{ bytes.Buffer }

func (bp *bufPrinter) Print(v ...interface{}) { fmt.Fprint(&bp.Buffer, append(v, "\n")...) }

func TestDefaultPrinter(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	stdlog.SetOutput(&buf)
	defer stdlog.SetOutput(os.Stderr)
	log := structlog.New(structlog.KeyUnit, "structlog")
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	pid := strconv.Itoa(os.Getpid())
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf structlog: `something happens` k1=v1 k2=v2 \t@ structlog_test.TestDefaultPrinter(output_test.go:25)\n"+
		"structlog.test["+pid+"] WRN structlog: `oops` \t@ structlog_test.TestDefaultPrinter(output_test.go:26)\n")
}

func TestPrinter(tt *testing.T) {
	t := check.T(tt)
	var buf bufPrinter
	log := structlog.New(structlog.KeyUnit, "structlog").SetPrinter(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	pid := strconv.Itoa(os.Getpid())
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf structlog: `something happens` k1=v1 k2=v2 \t@ structlog_test.TestPrinter(output_test.go:37)\n"+
		"structlog.test["+pid+"] WRN structlog: `oops` \t@ structlog_test.TestPrinter(output_test.go:38)\n")
}

func TestOutput(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	log := structlog.New(structlog.KeyUnit, "structlog").SetOutput(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	pid := strconv.Itoa(os.Getpid())
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf structlog: `something happens` k1=v1 k2=v2 \t@ structlog_test.TestOutput(output_test.go:49)\n"+
		"structlog.test["+pid+"] WRN structlog: `oops` \t@ structlog_test.TestOutput(output_test.go:50)\n")
}
