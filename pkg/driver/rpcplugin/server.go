package rpcplugin

import (
	"context"
	"encoding/json"
	pb "srpc/api/rpcdriver"
	"srpc/pkg/driver"
)

// newDriver 是实际协议驱动的构造函数（由插件主函数注入）
type newDriver func(json.RawMessage) (driver.Driver, error)

type server struct {
	pb.UnimplementedDriverServer

	create newDriver // 构造器
	drv    driver.Driver
}

// Init 由宿主进程先调用，把 param(json) 传进来。
// 这里才真正实例化具体协议驱动；可避免在插件启动阶段即做 IO。
func (s *server) Init(_ context.Context, in *pb.InitReq) (*pb.Ack, error) {
	if s.drv != nil {
		_ = s.drv.Close() // 二次 Init 先关掉旧实例
	}
	d, err := s.create(in.Param)
	if err != nil {
		return nil, err
	}
	s.drv = d
	return &pb.Ack{}, nil
}

func (s *server) Connect(ctx context.Context, _ *pb.Void) (*pb.Ack, error) {
	return &pb.Ack{}, s.drv.Connect(ctx)
}
func (s *server) Close(_ context.Context, _ *pb.Void) (*pb.Ack, error) {
	return &pb.Ack{}, s.drv.Close()
}
func (s *server) Health(_ context.Context, _ *pb.Void) (*pb.Ack, error) {
	return &pb.Ack{}, s.drv.Health()
}
func (s *server) Read(ctx context.Context, in *pb.ReadReq) (*pb.ReadResp, error) {
	pts, err := bytesToPoints(in.Points)
	if err != nil {
		return nil, err
	}
	vals, err := s.drv.Read(ctx, pts)
	if err != nil {
		return nil, err
	}
	return &pb.ReadResp{Values: valuesToBytes(vals)}, nil
}
func (s *server) Write(ctx context.Context, in *pb.WriteReq) (*pb.Ack, error) {
	var v driver.Value
	if err := json.Unmarshal(in.Value, &v); err != nil {
		return nil, err
	}
	return &pb.Ack{}, s.drv.Write(ctx, v)
}
