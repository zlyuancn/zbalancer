/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"errors"
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
	ins []interface{}
	mx  sync.RWMutex

	count      uint32
	cacheIndex []uint32 // 缓存的索引
}

func newWeightRandomBalancer() Balancer {
	return new(weightRandomBalancer)
}

func (b *weightRandomBalancer) Update(ins []interface{}, opt ...UpdateOption) {
	b.mx.Lock()
	defer b.mx.Unlock()

	// 选项
	opts := newUpdateOptions()
	opts.Apply(opt...)
	if len(opts.Weights) == 0 {
		opts.MakeDefaultWeight(len(ins))
	}
	if len(opts.Weights) != len(ins) {
		panic(errors.New("number of weights is inconsistent with the number of instances"))
	}

	// 计算权重
	var allWeight uint32
	b.ins = make([]interface{}, 0, len(ins))
	ends := make([]uint32, 0, len(ins)) // 实例在线段的结束位置列表
	for i, in := range ins {
		if opts.Weights[i] == 0 { // 权重为0忽略
			continue
		}
		allWeight += uint32(opts.Weights[i]) // 累加权重
		b.ins = append(b.ins, in)            // 每一个实例放在上一个实例的后面
		ends = append(ends, allWeight)       // 添加这个实例在线段的结束位置
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
		n := random.Int63()
		if isPowerOfTwo {
			score = uint32(n & int64(allWeight-1))
		} else {
			score = uint32(n % int64(allWeight))
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

func (b *weightRandomBalancer) Get(opt ...Option) (interface{}, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	l := len(b.ins)
	if l == 0 {
		return nil, NoInstanceErr
	}
	if l == 1 {
		return b.ins[0], nil
	}

	count := atomic.AddUint32(&b.count, 1)
	cacheIndex := count & (cacheIndexMode)
	return b.ins[b.cacheIndex[cacheIndex]], nil
}
