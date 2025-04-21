package main

import (
	hplugin "github.com/hashicorp/go-plugin"
	"srpc/internal/drivers/modbus"
	"srpc/pkg/driver/rpcplugin"
)

func main() {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: rpcplugin.Handshake,
		Plugins: map[string]hplugin.Plugin{
			"driver": &rpcplugin.DriverPlugin{
				Factory: modbus.NewDriver, // 注入具体协议驱动
			},
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
