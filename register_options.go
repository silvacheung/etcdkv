package etcdkv

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os"
	"runtime/debug"
	"time"
)

var registerErrorHandler = func(err error) {
	fmt.Fprintf(os.Stderr, "etcdkv register error:%v \n", err)
	debug.PrintStack()
}

func SetRegisterErrorHandler(fn func(error)) {
	if fn != nil {
		registerErrorHandler = fn
	}
}

type kvs struct {
	k string
	v string
}

// register注册器选项
type registerOption struct {
	client        *clientv3.Client
	namespace     string
	kvs           []kvs
	ttl           time.Duration
	leaseFaultTTL time.Duration
}

type RegisterOption func(*registerOption)

func RegisterClient(opts ...ClientOption) RegisterOption {
	client, err := newClientV3(opts...)
	if err != nil {
		registerErrorHandler(err)
	}
	return RegisterSetClient(client)
}

func RegisterSetClient(client *clientv3.Client) RegisterOption {
	return func(o *registerOption) {
		o.client = client
	}
}

func RegisterNamespace(namespace string) RegisterOption {
	return func(o *registerOption) {
		o.namespace = namespace
	}
}

func RegisterKvs(k, v string) RegisterOption {
	return func(o *registerOption) {
		o.kvs = append(o.kvs, kvs{k: k, v: v})
	}
}

func RegisterTTL(ttl time.Duration) RegisterOption {
	return func(o *registerOption) {
		o.ttl = ttl
	}
}

func RegisterLeaseFaultTTL(fault time.Duration) RegisterOption {
	return func(o *registerOption) {
		o.leaseFaultTTL = fault
	}
}
