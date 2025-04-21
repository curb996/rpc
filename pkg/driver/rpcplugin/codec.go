package rpcplugin

import (
	"encoding/json"

	"srpc/pkg/driver"
)

func marshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
func unmarshal[T any](b []byte) (out T, _ error) {
	err := json.Unmarshal(b, &out)
	return out, err
}

/* ---------- 对应 proto 字段转换 ---------- */

func pointsToBytes(pts []driver.Point) [][]byte {
	out := make([][]byte, len(pts))
	for i, p := range pts {
		out[i] = marshal(p)
	}
	return out
}
func bytesToPoints(bs [][]byte) ([]driver.Point, error) {
	out := make([]driver.Point, len(bs))
	for i, b := range bs {
		if err := json.Unmarshal(b, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}
func valuesToBytes(vs []driver.Value) [][]byte {
	out := make([][]byte, len(vs))
	for i, v := range vs {
		out[i] = marshal(v)
	}
	return out
}
func bytesToValues(bs [][]byte) ([]driver.Value, error) {
	out := make([]driver.Value, len(bs))
	for i, b := range bs {
		if err := json.Unmarshal(b, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}
