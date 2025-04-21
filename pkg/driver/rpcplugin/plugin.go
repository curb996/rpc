package rpcplugin

import (
	"context"
	"errors"

	hplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	pb "srpc/api/rpcdriver"
)

// Handshake 用于宿主 ↔ 插件 探测
var Handshake = hplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "DRIVER_PLUGIN",
	MagicCookieValue: "bess",
}

// ---------------- HashiCorp Plugin ----------------
type DriverPlugin struct {
	hplugin.Plugin
	Factory newDriver // 注入协议驱动的构造函数
}

// GRPCServer 在插件进程中执行，把 gRPC server 注册进来
func (p DriverPlugin) GRPCServer(b *hplugin.GRPCBroker, s *grpc.Server) error {
	if p.Factory == nil {
		return errors.New("DriverPlugin.Factory is nil")
	}
	pb.RegisterDriverServer(s, &server{create: p.Factory})
	return nil
}

// GRPCClient 在宿主进程中执行，构造一个 grpcClient
func (DriverPlugin) GRPCClient(
	_ context.Context, b *hplugin.GRPCBroker, cc *grpc.ClientConn,
) (interface{}, error) {
	return &grpcClient{c: pb.NewDriverClient(cc)}, nil
}
