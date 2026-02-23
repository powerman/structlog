package structlog_test

import (
	"errors"
	"io"
	"sync"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/structlog"
)

func TestGetErr(tt *testing.T) {
	t := check.T(tt)
	log := structlog.New().SetOutput(io.Discard)
	myerr := errors.New("my error") //nolint:err113 // By design.
	t.Err(log.Err(myerr), myerr)
	t.Err(log.Err(myerr, "err", io.EOF), myerr)
	t.Err(log.Err("fail", "err", io.EOF), io.EOF)
	t.Err(log.Err("fail", "a", 1, "myerr", myerr, "b", 2), myerr)
	t.Err(log.Err("fail", "a", 1, "b", 2), errors.New("fail")) //nolint:err113 // By design.
	t.Err(log.Err("fail", io.EOF, myerr), io.EOF)
	t.Err(log.Err("fail", io.EOF), io.EOF)
}

func TestNewNil(tt *testing.T) {
	t := check.T(tt)
	t.Panic(func() { (*structlog.Logger)(nil).New() }, "New called on nil *Logger")
}

// Just in case, not sure is it makes any sense to test this.
func TestRace1(_ *testing.T) {
	log := structlog.New().SetOutput(io.Discard).SetLogLevel(structlog.INF)
	log1 := log.New("key", "value")
	log2 := log.New()
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Go(func() { <-start; log.Err("failed") })
	wg.Go(func() { <-start; log.Warn("hmm") })
	wg.Go(func() { <-start; log1.Info("done") })
	wg.Go(func() { <-start; log2.Debug("dump") })
	close(start)
	wg.Wait()
}

// Just in case, not sure is it makes any sense to test this.
func TestRace2(_ *testing.T) {
	log0 := structlog.New().SetOutput(io.Discard)
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Go(func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 1); log.Err("failed") })
	wg.Go(func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 2); log.Warn("hmm") })
	wg.Go(func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 3); log.Info("done") })
	wg.Go(func() { <-start; log := log0.New(); log.SetDefaultKeyvals("key", 4); log.Debug("dump") })
	close(start)
	wg.Wait()
}
