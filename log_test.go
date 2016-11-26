package structlog_test

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"testing"

	"github.com/powerman/structlog"

	. "gopkg.in/check.v1"
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
	// c.Check(log.Err("fail", io.EOF, myerr), Equals, io.EOF)
}
