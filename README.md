# etcdkv
```
一个etcd的快速可拓展key/value事件监听包，可以快速用于服务注册与发现
也可以快速的实现一个简单的服务配置中心
```

# 例子
## 向etcd指定的命名空间下注册
```
	register := NewRegister(
		RegisterClient(
			ClientEndpoints("127.0.0.1:2379,127.0.0.1:2389,127.0.0.1:2399"),
			ClientDialKeepAliveTime(time.Second*5),
			ClientDialTimeout(time.Second*5),
			ClientDialKeepAliveTimeout(time.Second*5),
		),
		RegisterTTL(time.Second*10),
		RegisterNamespace("/"),
		RegisterKvs("1", "1111:1:1:1"),
		RegisterKvs("2", "2222:2:2:2"),
		RegisterKvs("3", "3333:3:3:3"),
	)

	defer register.Close()
	register.Start()
```
## 监听etcd指定的命名空间
```
	watcher := NewWatcher(
		WatcherClient(
			ClientEndpoints("127.0.0.1:2379,127.0.0.1:2389,127.0.0.1:2399"),
			ClientDialKeepAliveTime(time.Second*5),
			ClientDialTimeout(time.Second*5),
			ClientDialKeepAliveTimeout(time.Second*5),
		),
		WatcherNamespace("/"),
		WatcherTTL(time.Second*5),
		WatcherResolver(&PrintWatchKvResolver{}),
	)

	defer watcher.Close()
	watcher.Start()
```
## 自定义处理监听到变化的key/value信息
```
通过：etcdkv.WatcherResolver(&PrintWatchKvResolver{})这个设置可以指定监听到key/value变化时的处理方式
通过实现接口来实现自定义：

    type WatcherKvResolver interface {
    	Get(rawKv, namespace, key, value string, putTime, version int64) // 主动检测拉取时触发
    	Put(rawKv, namespace, key, value string, putTime, version int64) // 监听到有PUT事件时触发
    	Del(rawKv, namespace, key, value string, putTime, version int64) // 监听到有DELETE事件时触发
    }
```
## 设置注册和监听的错误处理器,可以自己处理错误
```
etcdkv.SetRegisterErrorHandler(func (err error) {
    // handle error
})

etcdkv.SetWatcherErrorHandler(func (err error) {
    // handle error
})
```

## 简单的测试
[注册](./register_test.go)
[监听](./watcher_test.go)
