package etcdkv

import (
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {

	watcher := NewWatcher(
		WatcherClient(
			ClientEndpoints("127.0.0.1:2379"),
			ClientDialKeepAliveTime(time.Second*5),
			ClientDialKeepAliveTimeout(time.Second*5),
		),
		WatcherNamespace(DefaultNamespace),
		WatcherTTL(time.Second*5),
		WatcherResolver(&PrintWatchKvResolver{}),
	)

	defer watcher.Close()
	watcher.Start()

	// 监听60秒退出
	time.Sleep(time.Second * 60)
}
