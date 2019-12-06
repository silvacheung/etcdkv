package etcdkv

import (
	"context"
	"crypto/tls"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"strings"
	"time"
)

const DefaultClientName = "globalClient"

// client客户端选项
type clientOption struct {
	name string
	cfg  clientv3.Config
}

type ClientOption func(*clientOption)

func ClientName(name string) ClientOption {
	return func(o *clientOption) {
		o.name = name
	}
}

func ClientEndpoints(endpoints string) ClientOption {
	return func(o *clientOption) {
		o.cfg.Endpoints = strings.Split(endpoints, ",")
	}
}

func ClientAutoSyncInterval(interval time.Duration) ClientOption {
	return func(o *clientOption) {
		o.cfg.AutoSyncInterval = interval
	}
}

func ClientDialTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOption) {
		o.cfg.DialTimeout = timeout
	}
}

func ClientDialKeepAliveTime(alive time.Duration) ClientOption {
	return func(o *clientOption) {
		o.cfg.DialKeepAliveTime = alive
	}
}

func ClientDialKeepAliveTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOption) {
		o.cfg.DialKeepAliveTimeout = timeout
	}
}

func ClientMaxCallSendMsgSize(size int) ClientOption {
	return func(o *clientOption) {
		o.cfg.MaxCallSendMsgSize = size
	}
}

func ClientMaxCallRecvMsgSize(size int) ClientOption {
	return func(o *clientOption) {
		o.cfg.MaxCallRecvMsgSize = size
	}
}

func ClientTLS(cfg *tls.Config) ClientOption {
	return func(o *clientOption) {
		o.cfg.TLS = cfg
	}
}

func ClientUsername(username string) ClientOption {
	return func(o *clientOption) {
		o.cfg.Username = username
	}
}

func ClientPassword(password string) ClientOption {
	return func(o *clientOption) {
		o.cfg.Password = password
	}
}

func ClientRejectOldCluster(reject bool) ClientOption {
	return func(o *clientOption) {
		o.cfg.RejectOldCluster = reject
	}
}

func ClientDialOptions(options ...grpc.DialOption) ClientOption {
	return func(o *clientOption) {
		o.cfg.DialOptions = options
	}
}

func ClientContext(ctx context.Context) ClientOption {
	return func(o *clientOption) {
		o.cfg.Context = ctx
	}
}
