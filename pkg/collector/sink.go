package collector

import (
	"context"
	"fmt"
	"srpc/pkg/driver"
)

// Sink 表示采集侧的输出目的地，可接入多种上行协议
type Sink interface {
	WriteBatch(ctx context.Context, vs []driver.Value) error
	Close() error
}

// ConsoleSink 用于 demo/测试
type ConsoleSink struct{}

func (ConsoleSink) WriteBatch(_ context.Context, vs []driver.Value) error {
	for _, v := range vs {
		// 打印缩略信息
		fmt.Printf("%s %-12s %v", v.Point.DeviceID, v.Point.Tag, v.Data)
	}
	return nil
}
func (ConsoleSink) Close() error { return nil }
