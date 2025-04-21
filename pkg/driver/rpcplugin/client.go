package rpcplugin

import (
	"context"
	"encoding/json"

	pb "srpc/api/rpcdriver"
	"srpc/pkg/driver"
)

type grpcClient struct {
	c pb.DriverClient
}

func (g *grpcClient) Init(p json.RawMessage) error {
	_, err := g.c.Init(context.Background(), &pb.InitReq{Param: p})
	return err
}
func (g *grpcClient) Connect(ctx context.Context) error {
	_, err := g.c.Connect(ctx, &pb.Void{})
	return err
}
func (g *grpcClient) Close() error {
	_, err := g.c.Close(context.Background(), &pb.Void{})
	return err
}
func (g *grpcClient) Health() error {
	_, err := g.c.Health(context.Background(), &pb.Void{})
	return err
}
func (g *grpcClient) Read(ctx context.Context, pts []driver.Point) ([]driver.Value, error) {
	reply, err := g.c.Read(ctx, &pb.ReadReq{Points: pointsToBytes(pts)})
	if err != nil {
		return nil, err
	}
	return bytesToValues(reply.Values)
}
func (g *grpcClient) Write(ctx context.Context, v driver.Value) error {
	b, _ := json.Marshal(v)
	_, err := g.c.Write(ctx, &pb.WriteReq{Value: b})
	return err
}
