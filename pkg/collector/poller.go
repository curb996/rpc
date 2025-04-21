package collector

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"srpc/pkg/driver"
)

//var log = zap.S()

type Poller struct {
	conf   DriverConf
	drv    driver.Driver
	sink   Sink
	wp     *Pool
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewPoller 自动实例化 driver 并绑定 sink
func NewPoller(cfg DriverConf, sink Sink) (*Poller, error) {
	drv, err := driver.Load(cfg.Type, cfg.Param)
	if err != nil {
		return nil, err
	}
	return &Poller{
		conf: cfg,
		drv:  drv,
		sink: sink,
		wp:   NewPool(64),
	}, nil
}

// Start —— blocking=false：后台运行
func (p *Poller) Start(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)
	interval := p.conf.Poll.Interval
	if interval == 0 {
		interval = time.Second
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		_ = p.drv.Connect(ctx)

		ticker := time.NewTicker(interval)
		span := otel.Tracer("poller").Start
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 用 worker pool 保障本轮超时不影响下一轮
				p.wp.Submit(func() {
					ctx2, cancel := context.WithTimeout(ctx, interval/2)
					defer cancel()
					var sCtx context.Context
					var sp trace.Span
					sCtx, sp = span(ctx2, "read")
					defer sp.End()

					values, err := p.drv.Read(sCtx, p.conf.Poll.Points)
					if err != nil {
						fmt.Printf("[%s] read err: %v", p.conf.ID, err)
						return
					}
					if err = p.sink.WriteBatch(sCtx, values); err != nil {
						fmt.Printf("sink err: %v", err)
					}
				})
			}
		}
	}()
}

// Stop 等待优雅关闭
func (p *Poller) Stop() {
	p.cancel()
	p.wp.Close()
	p.wg.Wait()
	_ = p.drv.Close()
}
