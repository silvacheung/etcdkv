package etcdkv

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"strconv"
	"sync"
	"time"
)

type Watcher struct {
	opt     *watcherOption
	ctx     context.Context
	cancel  context.CancelFunc
	watchCh clientv3.WatchChan
	wait    *sync.WaitGroup
	ticker  *time.Ticker
}

func NewWatcher(opts ...WatcherOption) *Watcher {

	opt := &watcherOption{}

	for _, optFun := range opts {
		optFun(opt)
	}

	if opt.client == nil {
		watcherErrorHandler(errors.New("etcdkv watcher client is empty"))
		return nil
	}

	if opt.ttl == 0 {
		opt.ttl = DefaultTTL
	}

	if opt.resolver == nil {
		opt.resolver = &PrintWatchKvResolver{}
	}

	wait, ticker := &sync.WaitGroup{}, time.NewTicker(opt.ttl)
	ctx, cancel := context.WithCancel(context.Background())
	watchCh := opt.client.Watch(ctx, opt.sepNamespace, clientv3.WithPrefix())

	return &Watcher{opt: opt, watchCh: watchCh, ctx: ctx, cancel: cancel, wait: wait, ticker: ticker}
}

func (w *Watcher) Start() {
	go w.start()
}

func (w *Watcher) Close() {
	w.close()
}

func (w *Watcher) start() {
	// First Get The Namespace
	w.fromNamespace()
	for {
		select {
		case change := <-w.watchCh: // 监听变化
			for _, event := range change.Events {
				kv := event.Kv
				if kv == nil {
					continue
				}
				k, v, putTime := w.parseKv(kv.Key, kv.Value)
				switch event.Type {
				case mvccpb.PUT:
					w.opt.resolver.Put(event.Kv.String(), w.opt.namespace, string(k), string(v), putTime, kv.Version)
				case mvccpb.DELETE:
					w.opt.resolver.Del(event.Kv.String(), w.opt.namespace, string(k), string(v), putTime, kv.Version)
				}
			}
		case <-w.ticker.C: // 定时拉取
			w.fromNamespace()
		case <-w.ctx.Done(): // 监听关闭
			log.Println("the watcher context is done")
			if err := w.opt.client.Close(); err != nil {
				watcherErrorHandler(err)
			}
			w.ticker.Stop()
			w.wait.Done()
			return
		}
	}
}

func (w *Watcher) close() {
	w.wait.Add(1)
	w.cancel()
	w.wait.Wait()
}

func (w *Watcher) fromNamespace() {
	if response, err := w.opt.client.Get(w.ctx, w.opt.sepNamespace, clientv3.WithPrefix()); err != nil {
		watcherErrorHandler(err)
	} else {
		for _, kv := range response.Kvs {
			k, v, putTime := w.parseKv(kv.Key, kv.Value)
			w.opt.resolver.Get(kv.String(), w.opt.namespace, string(k), string(v), putTime, kv.Version)
		}
	}
}

// 从key中去除namespace,从value中解析出该key的put的时间(s)和真实的值:key格式=/namespace/key;value格式=时间戳(s):value
func (w *Watcher) parseKv(key, value []byte) (k []byte, v []byte, putTime int64) {

	k, v = key, value

	// parse key
	if lenNamespace := len([]byte(w.opt.sepNamespace)); len(key) > lenNamespace {
		k = key[lenNamespace:]
	}

	// parse value and timestamp(s)
	if len(value) > 10 {
		if timestamp, err := strconv.Atoi(string(value[:10])); err == nil {
			v, putTime = value[11:], int64(timestamp)
		}
	}

	return
}
