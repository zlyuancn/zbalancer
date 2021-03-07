/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/7
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"sort"
	"sync"
)

type weightHashBalancer struct {
	allWeight uint32
	ins       []Instance
	ends      []uint32 //  实例在线段的结束位置列表
	hashFn    HashFn
	mx        sync.RWMutex
}

func newWeightHashBalancer() Balancer {
	return &weightHashBalancer{
		hashFn: DefaultHashFn,
	}
}

func (b *weightHashBalancer) Apply(opt ...BalancerOption) {
	opts := newBalancerOptions()
	opts.Apply(opt...)

	b.mx.Lock()
	b.hashFn = opts.HashFn
	b.mx.Unlock()
}

func (b *weightHashBalancer) Update(ins ...Instance) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.allWeight = 0
	b.ins = make([]Instance, 0, len(ins))
	b.ends = make([]uint32, 0, len(ins))
	for _, in := range ins {
		if in.Weight() == 0 { // 权重为0忽略
			continue
		}
		b.allWeight += uint32(in.Weight())   // 累加权重
		b.ins = append(b.ins, in)            // 每一个实例放在上一个实例的后面
		b.ends = append(b.ends, b.allWeight) // 添加这个实例在线段的结束位置
	}
}

// 二分搜索
func (b *weightHashBalancer) search(score uint32) int {
	return sort.Search(len(b.ends), func(i int) bool { return b.ends[i] > score })
}

func (b *weightHashBalancer) Get(opt ...Option) (Instance, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	l := len(b.ins)
	if l == 0 {
		return nil, NoInstanceErr
	}
	if l == 1 {
		return b.ins[0], nil
	}

	opts := newOptions()
	opts.Apply(opt...)

	hashValue := b.hashFn(opts.Key)

	var score uint32
	if b.allWeight&(b.allWeight-1) == 0 {
		score = hashValue & (b.allWeight - 1)
	} else {
		score = hashValue % b.allWeight
	}

	index := b.search(score)
	if index == l { // 环尾
		index = l - 1
	}
	return b.ins[index], nil
}
