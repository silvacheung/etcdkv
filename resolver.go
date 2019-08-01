package etcdkv

type WatcherKvResolver interface {
	Get(rawKv, namespace, key, value string, putTime, version int64)
	Put(rawKv, namespace, key, value string, putTime, version int64)
	Del(rawKv, namespace, key, value string, putTime, version int64)
}
