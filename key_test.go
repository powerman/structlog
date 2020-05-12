package structlog_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/powerman/check"
	"github.com/powerman/structlog"
)

func TestDefault(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf)
	log.Debug("defaults")
	t.Equal(buf.String(), "structlog.test["+pid+"] dbg "+unit+": `defaults` \t@ structlog_test.TestDefault(key_test.go:18)\n")
}

func TestDefaultJSON(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf).SetLogFormat(structlog.JSON)
	log.Debug("defaults")
	m := make(map[string]interface{})
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]interface{}{
		"_a": "structlog.test",
		"_f": "structlog_test.TestDefaultJSON",
		"_l": "dbg",
		"_m": "defaults",
		"_p": float64(os.Getpid()),
		"_s": "key_test.go:27",
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
	m := make(map[string]interface{})
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]interface{}{
		"_t": "02:04:05.123 UTC ",
		"_l": "ERR",
		"_m": "JSON",
	})
}
