package etcdkv

import "log"

// 默认实现接口,打印信息
type PrintWatchKvResolver struct{}

func (*PrintWatchKvResolver) Get(rawKv, namespace, key, value string, putTime, version int64) {
	log.Printf("the watcher kv resolver 'GET' event info: rawKv:%s, namespace:%s, key:%s, value:%s, putTime:%d, version:%d\n", rawKv, namespace, key, value, putTime, version)
}

func (*PrintWatchKvResolver) Put(rawKv, namespace, key, value string, putTime, version int64) {
	log.Printf("the watcher kv resolver 'PUT' event info: rawKv:%s, namespace:%s, key:%s, value:%s, putTime:%d, version:%d\n", rawKv, namespace, key, value, putTime, version)
}

func (*PrintWatchKvResolver) Del(rawKv, namespace, key, value string, putTime, version int64) {
	log.Printf("the watcher kv resolver 'DEL' event info: rawKv:%s, namespace:%s, key:%s, value:%s, putTime:%d, version:%d\n", rawKv, namespace, key, value, putTime, version)
}
