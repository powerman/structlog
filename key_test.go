package structlog_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/structlog"
)

func TestDefault(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf)
	_, err := os.Open("/no/such")
	log.Debug(err)
	log.SetDefaultKeyvals(structlog.KeyTime, structlog.Auto)
	log.Info("with time")
	t.Equal(buf.String(), ""+
		"structlog.test["+pid+"] dbg "+unit+": `open /no/such: no such file or directory` \t@ structlog_test.TestDefault(key_test.go:21)\n"+
		"Jan  2 03:04:05.123456 structlog.test["+pid+"] inf "+unit+": `with time` \t@ structlog_test.TestDefault(key_test.go:23)\n")
}

func TestDefaultJSON(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf).SetLogFormat(structlog.JSON)
	_, err := os.Open("/no/such")
	log.Debug(err)
	m := make(map[string]any)
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]any{
		"_a": "structlog.test",
		"_f": "structlog_test.TestDefaultJSON",
		"_l": "dbg",
		"_m": "open /no/such: no such file or directory",
		"_p": strconv.Itoa(os.Getpid()),
		"_s": "key_test.go:35",
		"_t": "Jan  2 02:04:05.123456",
		"_u": unit,
	})
}

func TestKeyTime(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.NewZeroLogger().SetOutput(&buf).SetPrefixKeys(structlog.KeyTime).
		SetDefaultKeyvals(structlog.KeyTime, structlog.Auto)
	log.Debug("zero")
	log.SetKeysFormat(map[string]string{
		structlog.KeyTime:    "%[2]s",
		structlog.KeyMessage: " %#[2]q",
	})
	log.Info("keys format")
	log.SetTimeFormat("15:04:05.999 MST ")
	log.Warn("time format")
	t.Equal(buf.String(), ""+
		" _t=Jan  2 03:04:05.123456 _m=zero\n"+
		"Jan  2 03:04:05.123456 `keys format`\n"+
		"03:04:05.123 CET  `time format`\n")

	buf.Reset()
	log.SetLogFormat(structlog.JSON)
	log.Err("JSON")
	m := make(map[string]any)
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]any{
		"_t": "02:04:05.123 UTC ",
		"_l": "ERR",
		"_m": "JSON",
	})
}
