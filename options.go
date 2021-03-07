/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

import (
	"hash/crc32"
)

// hash函数定义
type HashFn func(data []byte) uint32

// 更新选项
type updateOptions struct {
	Weights []uint8 // 权重, 权重值区间为0~255, 0表示永不使用
	HashFn  HashFn
}

// 更新选项定义
type UpdateOption func(opts *updateOptions)

func newUpdateOptions() *updateOptions {
	return &updateOptions{}
}

// 应用选项
func (opts *updateOptions) Apply(opt ...UpdateOption) {
	for _, o := range opt {
		o(opts)
	}
	if opts.HashFn == nil {
		opts.HashFn = crc32.ChecksumIEEE
	}
}

// 构建默认权重, 默认值为1
func (opts *updateOptions) MakeDefaultWeight(count int, def ...uint8) {
	opts.Weights = make([]uint8, count)
	for i := 0; i < count; i++ {
		if i < len(def) {
			opts.Weights[i] = def[i]
		} else {
			opts.Weights[i] = 1
		}
	}
}

// 设置权重
//
// 权重值区间为0~255, 0表示永不使用.
// 权重数量必须和实例数量一致, 默认所有实例权重为1.
//
// 仅以下平衡器生效:
//      WeightRandom
func WithUpdateWeights(weights []uint8) UpdateOption {
	return func(opts *updateOptions) {
		opts.Weights = weights
	}
}

// 设置hash函数
func WithUpdateHashFn(hashFn HashFn) UpdateOption {
	return func(opts *updateOptions) {
		opts.HashFn = hashFn
	}
}

// 获取选项
type options struct {
	Key []byte
}

// 获取选项定义
type Option func(opts *options)

func newOptions() *options {
	return &options{}
}

// 设置key
func WithKey(key []byte) Option {
	return func(opts *options) {
		opts.Key = key
	}
}

// 应用选项
func (opts *options) Apply(opt ...Option) {
	for _, o := range opt {
		o(opts)
	}
}
