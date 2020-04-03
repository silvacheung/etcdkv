package etcdkv

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"strings"
	"time"
)

var watcherErrorHandler = func(err error) {
	log.Printf("etcdkv watcher error:%v \n", err)
}

func SetWatcherErrorHandler(fn func(error)) {
	if fn != nil {
		watcherErrorHandler = fn
	}
}

// register注册器选项
type watcherOption struct {
	client       *clientv3.Client
	namespace    string
	sepNamespace string // like '/namespace/'
	ttl          time.Duration
	resolver     WatcherKvResolver
}

type WatcherOption func(*watcherOption)

func WatcherClient(opts ...ClientOption) WatcherOption {
	_, client, err := newClientV3(opts...)
	if err != nil {
		log.Println("the etcd watcher get client error:", err)
	}
	return WatcherSetClient(client)
}

func WatcherSetClient(client *clientv3.Client) WatcherOption {
	return func(o *watcherOption) {
		o.client = client
	}
}

func WatcherNamespace(namespace string) WatcherOption {
	return func(o *watcherOption) {
		o.namespace = namespace
		o.sepNamespace = namespaceWarp(namespace)
		if o.namespace == "" {
			o.namespace = o.sepNamespace
		}
	}
}

func WatcherTTL(ttl time.Duration) WatcherOption {
	return func(o *watcherOption) {
		o.ttl = ttl
	}
}

func WatcherResolver(resolver WatcherKvResolver) WatcherOption {
	return func(o *watcherOption) {
		o.resolver = resolver
	}
}

func namespaceWarp(namespace string) string {
	return fmt.Sprintf("%s/", strings.Trim(namespace, "/"))
}
