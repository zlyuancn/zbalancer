/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"sync"
	"sync/atomic"
)

type roundBalancer struct {
	incr   uint32 // 调用次数
	ins    []Instance
	target *targetSelector
	mx     sync.RWMutex
}

func newRoundBalancer() Balancer {
	return &roundBalancer{
		target: newTargetSelector(),
	}
}

func (b *roundBalancer) Apply(opt ...BalancerOption) {}

func (b *roundBalancer) Update(instances []Instance) {
	b.mx.Lock()
	b.incr = 0
	b.ins = instances
	b.target.Update(instances)
	b.mx.Unlock()
}

func (b *roundBalancer) Get(opt ...Option) (Instance, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	opts := newOptions()
	opts.Apply(opt...)

	if opts.Target != "" {
		return b.target.Get(opts.Target)
	}

	l := len(b.ins)
	if l == 0 {
		return nil, NoInstanceErr
	}
	if l == 1 {
		return b.ins[0], nil
	}

	count := atomic.AddUint32(&b.incr, 1) - 1
	var index int
	if l&(l-1) == 0 {
		index = int(count) & (l - 1)
	} else {
		index = int(count) % l
	}

	return b.ins[index], nil
}
