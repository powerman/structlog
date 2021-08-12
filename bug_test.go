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

func TestJSONMarshalError(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()
	var buf bytes.Buffer
	log := structlog.New().SetOutput(&buf).SetLogFormat(structlog.JSON)
	type V struct {
		S string
		F func()
		I int
	}
	v := V{S: "text", I: 42}

	log.Debug("msg", "v", v)
	m := make(map[string]interface{})
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]interface{}{
		"_a": "structlog.test",
		"_f": "structlog_test.TestJSONMarshalError",
		"_l": "dbg",
		"_m": "msg",
		"_p": strconv.Itoa(os.Getpid()),
		"_s": "bug_test.go:27",
		"_t": "Jan  2 02:04:05.123456",
		"_u": unit,
		"v":  "{text <nil> 42}",
	})

	buf.Reset()
	log.Debug(v)
	m = make(map[string]interface{})
	t.Nil(json.Unmarshal(buf.Bytes(), &m))
	t.DeepEqual(m, map[string]interface{}{
		"_a": "structlog.test",
		"_f": "structlog_test.TestJSONMarshalError",
		"_l": "dbg",
		"_m": "{text <nil> 42}",
		"_p": strconv.Itoa(os.Getpid()),
		"_s": "bug_test.go:43",
		"_t": "Jan  2 02:04:05.123456",
		"_u": unit,
	})
}
