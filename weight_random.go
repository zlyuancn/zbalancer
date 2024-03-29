/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	cacheIndexSize = 1 << 18 // 缓存索引大小
	cacheIndexMode = cacheIndexSize - 1
)

type weightRandomBalancer struct {
	ins    []Instance
	target *targetSelector
	mx     sync.RWMutex

	incr       uint32   // 调用次数
	cacheIndex []uint32 // 缓存的索引
}

func newWeightRandomBalancer() Balancer {
	return &weightRandomBalancer{
		target: newTargetSelector(),
	}
}

func (b *weightRandomBalancer) Apply(opt ...BalancerOption) {}

func (b *weightRandomBalancer) Update(instances []Instance) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.target.Update(instances)

	// 计算权重
	var allWeight uint32
	b.ins = make([]Instance, 0, len(instances))
	ends := make([]uint32, 0, len(instances)) // 实例在线段的结束位置列表
	for _, in := range instances {
		if in.Weight() == 0 { // 权重为0忽略
			continue
		}
		allWeight += uint32(in.Weight()) // 累加权重
		b.ins = append(b.ins, in)        // 每一个实例放在上一个实例的后面
		ends = append(ends, allWeight)   // 添加这个实例在线段的结束位置
	}

	// 实例总数小于2不需要计算缓存
	if len(b.ins) < 2 {
		return
	}

	// 缓存预先计算的index
	random := rand.New(rand.NewSource(time.Now().Unix())) // 随机生成器
	b.cacheIndex = make([]uint32, cacheIndexSize)
	isPowerOfTwo := allWeight&(allWeight-1) == 0 // 总权重是2的幂

	var score uint32
	for i := uint32(0); i < cacheIndexSize; i++ {
		// 生成随机数并计算分值
		n := random.Uint64()
		if isPowerOfTwo {
			score = uint32(n & uint64(allWeight-1))
		} else {
			score = uint32(n % uint64(allWeight))
		}

		// 根据分值搜索线段(实例)
		index := b.search(ends, score)
		if index == len(b.ins) { // 环尾, 虽然在这里不可能出现
			index = len(b.ins) - 1
		}
		b.cacheIndex[i] = uint32(index)
	}
}

// 二分搜索
func (b *weightRandomBalancer) search(scores []uint32, score uint32) int {
	return sort.Search(len(scores), func(i int) bool { return scores[i] > score })
}

func (b *weightRandomBalancer) Get(opt ...Option) (Instance, error) {
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

	count := atomic.AddUint32(&b.incr, 1)
	cacheIndex := count & (cacheIndexMode)
	return b.ins[b.cacheIndex[cacheIndex]], nil
}
