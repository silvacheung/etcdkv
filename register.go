package etcdkv

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/clientv3util"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	DefaultKey           = "default"
	DefaultValue         = "default"
	DefaultLeaseFaultTTL = time.Second * 30
	DefaultTTL           = time.Second * 60
)

type Register struct {
	opt     *registerOption
	ctx     context.Context
	cancel  context.CancelFunc
	wait    *sync.WaitGroup
	leaseId clientv3.LeaseID
	ticker  *time.Ticker
}

func NewRegister(opts ...RegisterOption) *Register {

	opt := &registerOption{}

	for _, optFun := range opts {
		optFun(opt)
	}

	if opt.client == nil {
		registerErrorHandler(errors.New("etcdkv register client is empty"))
		return nil
	}

	if len(opt.kvs) == 0 {
		opt.kvs = []kvs{{k: DefaultKey, v: DefaultValue}}
	}

	if opt.ttl == 0 {
		opt.ttl = DefaultTTL
	}

	if opt.leaseFaultTTL == 0 {
		opt.leaseFaultTTL = DefaultLeaseFaultTTL
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Register{opt: opt, ctx: ctx, cancel: cancel, wait: &sync.WaitGroup{}, ticker: time.NewTicker(opt.ttl)}
}

func (r *Register) Start() {
	go r.start()
}

func (r *Register) Close() {
	r.close()
}

func (r *Register) start() {
	for {
		r.register()
		select {
		case <-r.ticker.C: // 定时注册
		case <-r.ctx.Done(): // 监听关闭
			log.Println("etcdkv register context is done")
			r.unRegister()
			if err := r.opt.client.Close(); err != nil {
				registerErrorHandler(err)
			}
			r.ticker.Stop()
			r.wait.Done()
			return
		}
	}
}

func (r *Register) close() {
	r.wait.Add(1)
	r.cancel()
	r.wait.Wait()
}

// If Key Is Missing, Then Put Key
func (r *Register) register() {

	lease := clientv3.NewLease(r.opt.client)
	defer closeLease(lease)

	if leaseInfo, err := lease.TimeToLive(r.ctx, r.leaseId); err != nil {
		registerErrorHandler(err)
	} else if r.leaseId == 0 || leaseInfo.TTL == -1 { // 没有租约或者失效重新申请
		if leaseResponse, err := lease.Grant(r.ctx, int64((r.opt.ttl+r.opt.leaseFaultTTL)/time.Second)); err != nil {
			registerErrorHandler(err)
		} else {
			r.leaseId = leaseResponse.ID
		}
	} else { // 没有失效，延长租约
		if _, err := lease.KeepAlive(r.ctx, r.leaseId); err != nil {
			registerErrorHandler(err)
		}
	}

	if r.leaseId > 0 {
		kvClientV3, now := clientv3.NewKV(r.opt.client), time.Now().Unix()
		for _, kv := range r.opt.kvs {
			k, v := r.multiKv(now, kv.k, kv.v)
			if _, err := kvClientV3.Txn(r.ctx).If(clientv3util.KeyMissing(k)).
				Then(clientv3.OpPut(k, v, clientv3.WithLease(r.leaseId))).Commit(); err != nil {
				registerErrorHandler(err)
			}
		}
	}
}

// Revoke Lease
func (r *Register) unRegister() {
	lease := clientv3.NewLease(r.opt.client)
	defer closeLease(lease)
	if _, err := lease.Revoke(context.TODO(), r.leaseId); err != nil {
		registerErrorHandler(err)
	}
}

// return '/namespace/k, now(s):v'
func (r *Register) multiKv(now int64, k, v string) (key, value string) {
	key = fmt.Sprintf("%s/%s", strings.Trim(r.opt.namespace, "/"), k)
	value = fmt.Sprintf("%d:%s", now, v)
	return
}

func closeLease(lease clientv3.Lease) {
	if err := lease.Close(); err != nil {
		registerErrorHandler(err)
	}
}
