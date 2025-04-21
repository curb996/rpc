package collector

import (
	"encoding/json"
	"srpc/pkg/driver"
	"time"
)

// DriverConf 映射到 yaml -> drivers: [] …
type DriverConf struct {
	ID       string          `mapstructure:"id"`       // 唯一 ID
	Type     string          `mapstructure:"type"`     // 如 "modbus", "opcua"
	Param    json.RawMessage `mapstructure:"param"`    // 协议私有 JSON
	Poll     PollConf        `mapstructure:"poll"`     // 轮询策略
	Disabled bool            `mapstructure:"disabled"` // 支持配置停用
}

type PollConf struct {
	Interval time.Duration  `mapstructure:"interval"` // 采样周期
	Points   []driver.Point `mapstructure:"points"`   // 点表
}
