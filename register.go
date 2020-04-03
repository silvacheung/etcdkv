package etcdkv

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/clientv3util"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	DefaultKey   = "default"
	DefaultValue = "default"
)

type Register struct {
	opt     *registerOption
	ctx     context.Context
	cancel  context.CancelFunc
	wait    *sync.WaitGroup
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	keepCh  <-chan *clientv3.LeaseKeepAliveResponse
	closed  bool
}

func NewRegister(opts ...RegisterOption) *Register {

	opt := &registerOption{
		ttl: time.Second * 60,
	}

	for _, optFun := range opts {
		optFun(opt)
	}

	if opt.client == nil {
		registerErrorHandler(fmt.Errorf("client is empty"))
		return nil
	}

	if len(opt.kvs) == 0 {
		opt.kvs = []kvs{{k: DefaultKey, v: DefaultValue}}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Register{opt: opt, ctx: ctx, cancel: cancel, wait: &sync.WaitGroup{}}
}

func (r *Register) Start() {
	r.wait.Add(1)
	go func() {
		defer r.wait.Done()
		r.start()
	}()
}

func (r *Register) Close() {
	r.close()
}

func (r *Register) start() {
retry:
	if err := r.register(); err != nil {
		if !r.closed {
			err = fmt.Errorf("%v, retry again 3 second later", err)
			registerErrorHandler(err)
			time.Sleep(time.Second * 3)
			log.Printf("etcdkv regiter start retry ... \n")
			goto retry
		}
	}

	for {
		select {
		case alive, ok := <-r.keepCh:
			if (!ok || alive == nil) && !r.closed {
				err := fmt.Errorf("lease keep alive chan closed, retry again 3 second later")
				registerErrorHandler(err)
				time.Sleep(time.Second * 3)
				log.Printf("etcdkv register start retry ... \n")
				goto retry
			}
		case <-r.ctx.Done():
			log.Println("etcdkv register context is done")
			r.unRegister()
			r.closeLease()
			r.closeClient()
			return
		}
	}
}

func (r *Register) close() {
	r.closed = true
	r.cancel()
	r.wait.Wait()
}

// If Key Is Missing, Then Put Key
func (r *Register) register() error {

	r.lease = clientv3.NewLease(r.opt.client)

	if grant, err := r.lease.Grant(r.ctx, int64((r.opt.ttl)/time.Second)); err != nil {
		return err
	} else {
		r.leaseId = grant.ID
	}

	if r.leaseId > 0 {
		kvClientV3 := clientv3.NewKV(r.opt.client)
		for _, kv := range r.opt.kvs {
			k, v := r.multiKv(kv.k, kv.v)
			// 先删除已存在的KEY
			if _, err := kvClientV3.Txn(r.ctx).
				If(clientv3util.KeyExists(k)).
				Then(clientv3.OpDelete(k)).Commit(); err != nil {
				return err
			}
			// 注册新的KEY和租约绑定
			if txn, err := kvClientV3.Txn(r.ctx).
				If(clientv3util.KeyMissing(k)).
				Then(clientv3.OpPut(k, v, clientv3.WithLease(r.leaseId))).Commit(); err != nil {
				return err
			} else {
				log.Printf("etcdkv register opput key(%s) value(%s) succeeded(%v) \n", kv.k, kv.v, txn.Succeeded)
			}
		}
	}

	if keepCh, err := r.lease.KeepAlive(r.ctx, r.leaseId); err != nil {
		return err
	} else {
		r.keepCh = keepCh
		return nil
	}
}

// Revoke Lease
func (r *Register) unRegister() {
	if _, err := r.lease.Revoke(context.TODO(), r.leaseId); err != nil {
		registerErrorHandler(err)
	}
}

// return '/namespace/k, v'
func (r *Register) multiKv(k, v string) (key, value string) {
	key = fmt.Sprintf("%s/%s", strings.Trim(r.opt.namespace, "/"), k)
	return key, v
}

func (r *Register) closeLease() {
	if err := r.lease.Close(); err != nil {
		registerErrorHandler(err)
	}
}

func (r *Register) closeClient() {
	if err := r.opt.client.Close(); err != nil {
		registerErrorHandler(err)
	}
}
