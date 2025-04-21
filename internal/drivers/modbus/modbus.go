package modbus

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"srpc/pkg/driver"
	"time"

	"github.com/goburrow/modbus"
)

// Config 私有字段，根据需要扩展
type Config struct {
	Mode       string `json:"mode"`       // rtu or tcp
	SerialPort string `json:"serialPort"` // e.g. /dev/ttyUSB0
	BaudRate   int    `json:"baudRate"`   // 9600
	SlaveID    byte   `json:"slaveId"`    // 1
	TCPAddr    string `json:"tcpAddr"`    // 192.168.1.10:502
	TimeoutMs  int    `json:"timeoutMs"`  // 1000
}

// driverImpl 满足 driver.Driver
type driverImpl struct {
	cfg Config
	cl  modbus.Client
}

// NewDriver 插件导出符号(必须)
func NewDriver(raw json.RawMessage) (driver.Driver, error) {
	var c Config
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return &driverImpl{cfg: c}, nil
}

// Connect 根据 cfg.Mode 建立底层链接
func (d *driverImpl) Connect(_ context.Context) error {
	var handler modbus.ClientHandler
	if d.cfg.Mode == "rtu" {
		h := modbus.NewRTUClientHandler(d.cfg.SerialPort)
		h.BaudRate = d.cfg.BaudRate
		h.SlaveId = d.cfg.SlaveID
		h.Timeout = time.Duration(d.cfg.TimeoutMs) * time.Millisecond
		handler = h
	} else {
		h := modbus.NewTCPClientHandler(d.cfg.TCPAddr)
		h.Timeout = time.Duration(d.cfg.TimeoutMs) * time.Millisecond
		h.SlaveId = d.cfg.SlaveID
		handler = h
	}

	//if err := handler.Connect(); err != nil {
	//	return err
	//}
	d.cl = modbus.NewClient(handler)

	return nil
}

func (d *driverImpl) Close() error {
	if h, ok := d.cl.(interface{ Close() error }); ok {
		return h.Close()
	}
	return nil
}

func (d *driverImpl) Health() error { return nil }

func (d *driverImpl) Read(ctx context.Context, pts []driver.Point) ([]driver.Value, error) {
	// 简易版：逐点读取（可自行批量优化）
	out := make([]driver.Value, 0, len(pts))
	for _, p := range pts {
		var reg []byte
		var err error
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		addr := p.Address
		switch p.Type {
		case driver.Float32, driver.Int16, driver.Uint16:
			reg, err = d.cl.ReadHoldingRegisters(uint16(addr), 1)
		case driver.Int32, driver.Uint32, driver.Float64:
			reg, err = d.cl.ReadHoldingRegisters(uint16(addr), 2)
		default:
			err = fmt.Errorf("unsupport type %v", p.Type)
		}
		if err != nil {
			return nil, err
		}
		v := parseRegister(p, reg)
		out = append(out, driver.Value{
			Point: p,
			Ts:    time.Now(),
			Data:  v,
		})
	}
	return out, nil
}

func (d *driverImpl) Write(ctx context.Context, v driver.Value) error {
	return driver.ErrUnsupported
}

// parseRegister 将寄存器raw -> any
func parseRegister(p driver.Point, reg []byte) any {
	var val any
	switch p.Type {
	case driver.Float32:
		val = float32(binary.BigEndian.Uint16(reg)) * float32(p.Scale)
		// ... 省略其他类型
	}
	return val
}
