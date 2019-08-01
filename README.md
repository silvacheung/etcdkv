# etcdkv
```
一个etcd的快速可拓展key/value事件监听包，可以快速用于服务注册与发现
也可以快速的实现一个简单的服务配置中心
```
# 依赖
```
使用命令快速解决：go mod tidy 
或者查看go.mod自己解决依赖包
```
# 例子
## 向etcd指定的命名空间下注册
```
    register := etcdkv.NewRegister(
		etcdkv.RegisterClient( // 设置etcd客户端
			etcdkv.ClientEndpoints("127.0.0.1:2379"),
			etcdkv.ClientDialKeepAliveTime(time.Second*5),
			etcdkv.ClientDialKeepAliveTimeout(time.Second*5),
		),
		etcdkv.RegisterTTL(time.Second*50), // 设置TTL,主要用于租约
		etcdkv.RegisterLeaseFaultTTL(time.Second*5), // 设置一个租约的容错时间,即延长租约失效时间
		etcdkv.RegisterNamespace(etcdkv.DefaultNamespace), // 设置注册的命名空间
		etcdkv.RegisterKvs("1", "1111:1:1:1"), // 要注册的一系列key/value,支持多个
		etcdkv.RegisterKvs("2", "2222:2:2:2"),
		etcdkv.RegisterKvs("3", "3333:3:3:3"),
	)

	defer register.Close()
	register.Start()
```
## 监听etcd指定的命名空间
```
    watcher := etcdkv.NewWatcher(
		etcdkv.WatcherClient( // 设置etcd客户端
			etcdkv.ClientEndpoints("127.0.0.1:2379"),
			etcdkv.ClientDialKeepAliveTime(time.Second*5),
			etcdkv.ClientDialKeepAliveTimeout(time.Second*5),
		),
		etcdkv.WatcherNamespace(DefaultNamespace), // 设置监听的命名空间
		etcdkv.WatcherTTL(time.Second*5), // 设置主动检测时间,不设置则不自动检测
		etcdkv.WatcherResolver(&PrintWatchKvResolver{}), // 设置监听事件时的key/value处理器
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