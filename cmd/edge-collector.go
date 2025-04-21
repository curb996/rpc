package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"srpc/pkg/collector"
)

func initLogger() {
	l, _ := zap.NewProduction()
	zap.ReplaceGlobals(l)
}

func main() {
	initLogger()

	// --------------------------------------------------
	viper.SetConfigName("edge") // edge.yaml
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatal(err)
	}

	var cfg struct {
		Drivers []collector.DriverConf `mapstructure:"drivers"`
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		zap.S().Fatal(err)
	}
	// --------------------------------------------------
	sink := collector.ConsoleSink{} // 可替换 NatsSink / MqttSink

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var pollers []*collector.Poller
	for _, d := range cfg.Drivers {
		if d.Disabled {
			continue
		}
		p, err := collector.NewPoller(d, sink)
		if err != nil {
			zap.S().Errorf("poller err: %v", err)
			continue
		}
		p.Start(ctx)
		pollers = append(pollers, p)
		zap.S().Infof("poller %s start OK", d.ID)
	}

	<-ctx.Done()
	zap.S().Info("shutting down ...")
	for _, p := range pollers {
		p.Stop()
	}
	time.Sleep(100 * time.Millisecond)
}
