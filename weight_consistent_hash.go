/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/7
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

type weightConsistentHashBalancer struct {
	incr    uint32 // 调用次数
	ends    []uint32
	target  *targetSelector
	hashMap map[uint32]Instance
	hashFn  HashFn
	mx      sync.RWMutex
}

func newWeightConsistentHashBalancer() Balancer {
	return &weightConsistentHashBalancer{
		hashFn: DefaultHashFn,
		target: newTargetSelector(),
	}
}

func (b *weightConsistentHashBalancer) Apply(opt ...BalancerOption) {
	opts := newBalancerOptions()
	opts.Apply(opt...)

	b.mx.Lock()
	b.hashFn = opts.HashFn
	b.mx.Unlock()
}

func (b *weightConsistentHashBalancer) Update(instances []Instance) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.ends = make([]uint32, 0)
	b.target.Update(instances)
	b.hashMap = make(map[uint32]Instance)
	for _, in := range instances {
		if in.Name() == "" {
			panic(errors.New("instance name is empty"))
		}
		if in.Weight() == 0 { // 权重为0忽略
			continue
		}
		for shard := 0; shard < int(in.Weight()); shard++ {
			hashValue := b.hashFn([]byte(fmt.Sprintf("%s_%d", in.Name(), shard))) // 计算分片的hash值
			b.ends = append(b.ends, hashValue)                                    // 这个hash值将作为环上的一点
			b.hashMap[hashValue] = in                                             // 保存hash值指向的实例
		}
	}

	if len(b.ends) > 0 {
		sort.Slice(b.ends, func(i, j int) bool {
			return b.ends[i] < b.ends[j]
		})
	}
}

// 二分搜索
func (b *weightConsistentHashBalancer) search(score uint32) int {
	return sort.Search(len(b.ends), func(i int) bool { return b.ends[i] >= score })
}

func (b *weightConsistentHashBalancer) Get(opt ...Option) (Instance, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	opts := newOptions()
	opts.Apply(opt...)

	if opts.Target != "" {
		return b.target.Get(opts.Target)
	}

	if len(b.hashMap) == 0 {
		return nil, NoInstanceErr
	}

	if len(opts.HashKey) == 0 {
		incr := atomic.AddUint32(&b.incr, 1)
		opts.HashKey = []byte(strconv.Itoa(int(incr)))
	}

	hashValue := b.hashFn(opts.HashKey)
	endIndex := b.search(hashValue)
	if endIndex == len(b.ends) { // 环尾
		endIndex = 0
	}
	return b.hashMap[b.ends[endIndex]], nil
}
