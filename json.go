package structlog

import (
	"encoding/json"
	"fmt"
)

type kvs map[string]interface{}

func (kv kvs) MarshalJSON() ([]byte, error) {
	safe := make(map[string]string, len(kv))
	for k, v := range kv {
		safe[k] = fmt.Sprint(v)
	}
	return json.Marshal(safe)
}
