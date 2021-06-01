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
	count uint32
	ins   []Instance
	mx    sync.RWMutex
}

func newRoundBalancer() Balancer {
	return new(roundBalancer)
}

func (b *roundBalancer) Apply(opt ...BalancerOption) {}

func (b *roundBalancer) Update(instances []Instance) {
	b.mx.Lock()
	b.ins = instances
	b.mx.Unlock()
}

func (b *roundBalancer) Get(opt ...Option) (Instance, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	l := len(b.ins)
	if l == 0 {
		return nil, NoInstanceErr
	}
	if l == 1 {
		return b.ins[0], nil
	}

	count := atomic.AddUint32(&b.count, 1) - 1
	var index int
	if l&(l-1) == 0 {
		index = int(count) & (l - 1)
	} else {
		index = int(count) % l
	}

	return b.ins[index], nil
}
