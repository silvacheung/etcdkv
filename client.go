package etcdkv

import (
	"github.com/coreos/etcd/clientv3"
)

var clientV3 *clientv3.Client

func newClientV3(opts ...ClientOption) (*clientv3.Client, error) {
	clientOpt := &clientOption{}
	for _, opt := range opts {
		opt(clientOpt)
	}
	return clientv3.New(clientOpt.cfg)
}

func NewClientV3(opts ...ClientOption) (*clientv3.Client, error) {
	var err error
	clientV3, err = newClientV3(opts...)
	return clientV3, err
}

func ClientV3() (*clientv3.Client, bool) {
	var ok bool
	if clientV3 != nil {
		ok = true
	}
	return clientV3, ok
}
