package driver

import (
	"context"
	"time"
)

// PointType 适配工控常见数据类型
type PointType uint8

const (
	Bool PointType = iota + 1
	Int16
	Int32
	Uint16
	Uint32
	Float32
	Float64
	String
)

// Point 描述一条采集点
type Point struct {
	DeviceID string    // 逻辑设备 ID
	Tag      string    // 外部可读名称
	Address  uint32    // 寄存器/节点地址
	Type     PointType // 数据类型
	Scale    float64   // 修正倍数（=1 表示无缩放）
}

// Value 封装一条实时值
type Value struct {
	Point
	Ts   time.Time
	Data any // 根据 Point.Type 再做断言
}

// Driver 所有协议驱动必须实现的接口
type Driver interface {
	// Connect/Close 生命周期，在 poller 启动或重连时调用
	Connect(ctx context.Context) error
	Close() error

	// Read 批量读取；同协议下可自行做管脚合并/批量优化
	Read(ctx context.Context, points []Point) ([]Value, error)

	// Write 下置单点；如协议不支持可返回 driver.ErrUnsupported
	Write(ctx context.Context, v Value) error

	// Health 用于探活（不阻塞 IO）
	Health() error
}
