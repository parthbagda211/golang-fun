package backend

import (
	"context"
	"go.uber.org/zap"
	"loadbalancer-go/Utils"
	_ "loadbalancer-go/Utils"
	"net"
	"net/url"
)

func IsBackendAlive(ctx context.Context, aliveChannel chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		Utils.Logger.Error("Backend is not alive", zap.String("url", u.String()), zap.Error(err))
		aliveChannel <- false
		return
	}
	_ = conn.Close()
	aliveChannel <- true
}
