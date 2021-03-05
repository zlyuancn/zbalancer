/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

type updateOptions struct {
	Weights []uint8 // 权重, 权重值区间为0~255, 0表示永不使用
}

type UpdateOption func(opts *updateOptions)

func newUpdateOptions() *updateOptions {
	return &updateOptions{}
}

// 构建默认权重, 默认值为1
func (opts *updateOptions) makeDefaultWeight(count int, def ...uint8) {
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

type options struct {
}

type Option func(opts *options)

func newOptions() *options {
	return &options{}
}
