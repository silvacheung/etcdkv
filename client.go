package etcdkv

import (
	"github.com/coreos/etcd/clientv3"
)

var clientV3 *clientv3.Client

func NewClientV3(opts ...ClientOption) (*clientv3.Client, error) {
	clientOpt := &clientOption{}
	for _, opt := range opts {
		opt(clientOpt)
	}
	var err error
	clientV3, err = clientv3.New(clientOpt.cfg)
	return clientV3, err
}

func ClientV3() (*clientv3.Client, bool) {
	var ok bool
	if clientV3 != nil {
		ok = true
	}
	return clientV3, ok
}
