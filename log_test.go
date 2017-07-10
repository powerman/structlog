package structlog_test

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"sync"
	"testing"

	"github.com/LiflandGaming/structlog"

	"bytes"
	. "gopkg.in/check.v1"
	"regexp"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *C) {
	log.SetOutput(ioutil.Discard)
}

func (s *TestSuite) SetUpTest(c *C) {}

func (s *TestSuite) TearDownTest(c *C) {}

func (s *TestSuite) TearDownSuite(c *C) {}

func (s *TestSuite) TestGetErr(c *C) {
	log := structlog.New()
	myerr := errors.New("my error")
	c.Check(log.Err(myerr), Equals, myerr)
	c.Check(log.Err(myerr, "err", io.EOF), Equals, myerr)
	c.Check(log.Err("fail", "err", io.EOF), Equals, io.EOF)
	c.Check(log.Err("fail", "a", 1, "myerr", myerr, "b", 2), Equals, myerr)
	c.Check(log.Err("fail", "a", 1, "b", 2), DeepEquals, errors.New("fail"))
	c.Check(log.Err("fail", io.EOF, myerr), Equals, io.EOF)
	c.Check(log.Err("fail", io.EOF), Equals, io.EOF)
}

func (s *TestSuite) TestNewNil(c *C) {
	preamble := func() (log *structlog.Logger) {
		log = log.New()
		log.Err("log without parent")
		return
	}
	c.Check(preamble, Panics, "New called on nil *Logger")
}

// Just in case, not sure is it makes any sense to test this.
func (s *TestSuite) TestRace1(c *C) {
	structlog.DefaultLogger.SetLogLevel(structlog.INF)
	log := structlog.New()
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
	c.Succeed()
}

// Just in case, not sure is it makes any sense to test this.
func (s *TestSuite) TestRace2(c *C) {
	var wg sync.WaitGroup
	wg.Add(4)
	start := make(chan struct{})
	go func() {
		<-start
		log := structlog.New()
		log.SetDefaultKeyvals("key", 1)
		log.Error("failed")
		wg.Done()
	}()
	go func() { <-start; log := structlog.New(); log.SetDefaultKeyvals("key", 2); log.Warn("hmm"); wg.Done() }()
	go func() { <-start; log := structlog.New(); log.SetDefaultKeyvals("key", 3); log.Info("done"); wg.Done() }()
	go func() { <-start; log := structlog.New(); log.SetDefaultKeyvals("key", 4); log.Debug("dump"); wg.Done() }()
	close(start)
	wg.Wait()
	c.Succeed()
}

func (s *TestSuite) TestWithFields(c *C) {
	buf := make([]byte, 200)
	writer := bytes.NewBuffer(buf)
	log.SetOutput(writer)
	defer log.SetOutput(ioutil.Discard)

	log := structlog.New()
	log.SetLogFormat(structlog.JSONFormatte)
	log1 := log.WithField("key", "value")
	log2 := log1.WithField("key2", "value2")

	log2.Printf("Test %v", "message")

	// workaround, because 	c.Assert(writer.String(), Matches, `.*key.*`) not working
	ok, err := regexp.MatchString(`.*"key":"value".*`, writer.String())
	if err != nil {
		c.Error(err)
	}
	c.Assert(ok, Equals, true)

	ok, err = regexp.MatchString(`.*"key2":"value2".*`, writer.String())
	if err != nil {
		c.Error(err)
	}
	c.Assert(ok, Equals, true)

	ok, err = regexp.MatchString(`.*"Test message".*`, writer.String())
	if err != nil {
		c.Error(err)
	}
	c.Assert(ok, Equals, true)

}
