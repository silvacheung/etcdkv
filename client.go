package etcdkv

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
)

var clientV3Map sync.Map

func newClientV3(opts ...ClientOption) (string, *clientv3.Client, error) {
	clientOpt := &clientOption{name: DefaultClientName}
	for _, opt := range opts {
		opt(clientOpt)
	}
	c, err := clientv3.New(clientOpt.cfg)
	return clientOpt.name, c, err
}

func NewClientV3(opts ...ClientOption) error {
	if name, clientV3, err := newClientV3(opts...); err != nil {
		return err
	} else {
		clientV3Map.Store(name, clientV3)
	}
	return nil
}

func ClientV3(name string) *clientv3.Client {
	if c, ok := clientV3Map.Load(name); ok {
		return c.(*clientv3.Client)
	} else {
		return nil
	}
}
