package structlog_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/structlog"
)

type nilPrinter struct{}

func (np nilPrinter) Print(v ...interface{}) {}

type bufPrinter struct{ bytes.Buffer }

func (bp *bufPrinter) Print(v ...interface{}) { fmt.Fprint(&bp.Buffer, append(v, "\n")...) }

func TestDefaultPrinter(tt *testing.T) {
	t := check.T(tt)
	var buf bytes.Buffer
	stdlog.SetOutput(&buf)
	defer stdlog.SetOutput(os.Stderr)
	log := structlog.New()
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	pid := strconv.Itoa(os.Getpid())
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf structlog: `something happens` k1=v1 k2=v2 \t@ structlog_test.TestDefaultPrinter(log_test.go:32)\n"+
		"structlog.test["+pid+"] WRN structlog: `oops` \t@ structlog_test.TestDefaultPrinter(log_test.go:33)\n")
}

func TestPrinter(tt *testing.T) {
	t := check.T(tt)
	var buf bufPrinter
	log := structlog.New().SetPrinter(&buf)
	log.Info("something happens", "k1", "v1", "k2", "v2")
	log.Warn("oops")
	pid := strconv.Itoa(os.Getpid())
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] inf structlog: `something happens` k1=v1 k2=v2 \t@ structlog_test.TestPrinter(log_test.go:44)\n"+
		"structlog.test["+pid+"] WRN structlog: `oops` \t@ structlog_test.TestPrinter(log_test.go:45)\n")
}

func TestGetErr(tt *testing.T) {
	t := check.T(tt)
	log := structlog.New().SetPrinter(nilPrinter{})
	myerr := errors.New("my error")
	t.Err(log.Err(myerr), myerr)
	t.Err(log.Err(myerr, "err", io.EOF), myerr)
	t.Err(log.Err("fail", "err", io.EOF), io.EOF)
	t.Err(log.Err("fail", "a", 1, "myerr", myerr, "b", 2), myerr)
	t.Err(log.Err("fail", "a", 1, "b", 2), errors.New("fail"))
	t.Err(log.Err("fail", io.EOF, myerr), io.EOF)
	t.Err(log.Err("fail", io.EOF), io.EOF)
}

func TestNewNil(tt *testing.T) {
	t := check.T(tt)
	t.Panic(func() { (*structlog.Logger)(nil).New() }, "New called on nil *Logger")
}

// Just in case, not sure is it makes any sense to test this.
func TestRace1(t *testing.T) {
	log := structlog.New().SetPrinter(nilPrinter{}).SetLogLevel(structlog.INF)
	log1 := log.New("key", "value")
	log2 := log.New()
	var wg sync.WaitGroup
	wg.Add(4)
	start := make(chan struct{})
	go func() { <-start; log.Err("failed"); wg.Done() }()
	go func() { <-start; log.Warn("hmm"); wg.Done() }()
	go func() { <-start; log1.Info("done"); wg.Done() }()
	go func() { <-start; log2.Debug("dump"); wg.Done() }()
	close(start)
	wg.Wait()
}

// Just in case, not sure is it makes any sense to test this.
func TestRace2(t *testing.T) {
	log0 := structlog.New().SetPrinter(nilPrinter{})
	var wg sync.WaitGroup
	wg.Add(4)
	start := make(chan struct{})
	go func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 1); log.Err("failed"); wg.Done() }()
	go func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 2); log.Warn("hmm"); wg.Done() }()
	go func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 3); log.Info("done"); wg.Done() }()
	go func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 4); log.Debug("dump"); wg.Done() }()
	close(start)
	wg.Wait()
}
