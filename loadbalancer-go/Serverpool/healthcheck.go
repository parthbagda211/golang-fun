package Serverpool

import (
	"context"
	"loadbalancer-go/Utils"
	_ "loadbalancer-go/Utils"
	"time"
)

func HealthCheck(ctx context.Context, sp ServerPool) {
	t := time.NewTicker(time.Second * 20)
	Utils.Logger.Info("Starting health check")
	for {
		select {
		case <-t.C:
			go HealthCheck(ctx, sp)
		case <-ctx.Done():
			Utils.Logger.Info("Closing Health Check")
			return
		}
	}
}
