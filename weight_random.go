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
	"time"
)

type weightRandomBalancer struct {
	allWeight uint32
	ins       []interface{}
	scores    []uint32
	mx        sync.Mutex
	random    *rand.Rand
}

func newWeightRandomBalancer() Balancer {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	return &weightRandomBalancer{
		random: random,
	}
}

func (b *weightRandomBalancer) Update(ins []interface{}, opt ...UpdateOption) {
	b.mx.Lock()
	defer b.mx.Unlock()

	opts := newUpdateOptions()
	opts.Apply(opt...)
	if len(opts.Weights) == 0 {
		opts.MakeDefaultWeight(len(ins))
	}
	if len(opts.Weights) != len(ins) {
		panic(errors.New("number of weights is inconsistent with the number of instances"))
	}

	b.allWeight = 0
	b.ins = make([]interface{}, 0, len(ins))
	b.scores = make([]uint32, 0, len(ins))
	for i, in := range ins {
		if opts.Weights[i] == 0 {
			continue
		}
		b.ins = append(b.ins, in)
		b.allWeight += uint32(opts.Weights[i])   // 这个值是记录所有权重的累加
		b.scores = append(b.scores, b.allWeight) // 每一个实例放在上一个实例的后面
	}
}

// 二分搜索
func (b *weightRandomBalancer) search(score uint32) int {
	return sort.Search(len(b.scores), func(i int) bool { return b.scores[i] > score })
}

func (b *weightRandomBalancer) Get(opt ...Option) (interface{}, error) {
	b.mx.Lock()
	defer b.mx.Unlock()

	l := len(b.ins)
	if l == 0 {
		return nil, NoInstanceErr
	}
	if l == 1 {
		return b.ins[0], nil
	}

	score := b.random.Int31n(int32(b.allWeight))
	index := b.search(uint32(score))
	if index == l { // 环尾
		index = l - 1
	}
	return b.ins[index], nil
}
