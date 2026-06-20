package handlers

import (
	"encoding/json"
)

// rawToInterface 把 json.RawMessage 安全转为 interface{}，空则返回 nil。
func rawToInterface(raw json.RawMessage) interface{} {
	if len(raw) == 0 {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil
	}
	return v
}
