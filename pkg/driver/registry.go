package driver

import (
	"encoding/json"
	"fmt"
	"plugin"
	"sync"
)

// newDriverSymbol 是插件内必须导出的符号名
const newDriverSymbol = "NewDriver"

// factory 在 registry 中保存构造函数
type factory func(json.RawMessage) (Driver, error)

var (
	mu       sync.RWMutex
	registry = map[string]factory{}
)

// Register 供静态编译场景调用。例如内置 OPC‑UA 驱动
func Register(name string, f factory) { // nolint: revive
	mu.Lock()
	registry[name] = f
	mu.Unlock()
}

// Load 根据 name 建立实例。若本进程未静态注册，则尝试 .so 插件
// 例如 name=modbus → ./plugins/modbus.so
func Load(name string, raw json.RawMessage) (Driver, error) {
	mu.RLock()
	if f, ok := registry[name]; ok {
		mu.RUnlock()
		return f(raw)
	}
	mu.RUnlock()

	// 动态插件
	p, err := plugin.Open(fmt.Sprintf("%s.so", name))
	if err != nil {
		return nil, fmt.Errorf("open plugin: %w", err)
	}
	sym, err := p.Lookup(newDriverSymbol)
	if err != nil {
		return nil, fmt.Errorf("lookup %s: %w", newDriverSymbol, err)
	}
	creator, ok := sym.(func(json.RawMessage) (Driver, error))
	if !ok {
		return nil, fmt.Errorf("symbol %s has wrong signature", newDriverSymbol)
	}
	return creator(raw)
}
