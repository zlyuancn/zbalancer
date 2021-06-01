package zbalancer

import (
	"hash/crc32"
)

// hash函数定义
type HashFn func(data []byte) uint32

// 默认hash函数
var DefaultHashFn = crc32.ChecksumIEEE

// 平衡器选项
type balancerOptions struct {
	HashFn HashFn
}

// 平衡器选项定义
type BalancerOption func(opts *balancerOptions)

func newBalancerOptions() *balancerOptions {
	return &balancerOptions{}
}

// 应用选项
func (opts *balancerOptions) Apply(opt ...BalancerOption) {
	for _, o := range opt {
		o(opts)
	}
	if opts.HashFn == nil {
		opts.HashFn = DefaultHashFn
	}
}

// 设置hash函数
func WithBalancerHashFn(hashFn HashFn) BalancerOption {
	return func(opts *balancerOptions) {
		opts.HashFn = hashFn
	}
}
