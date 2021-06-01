/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

// 获取选项
type getOptions struct {
	Key []byte
}

// 获取选项定义
type Option func(opts *getOptions)

func newOptions() *getOptions {
	return &getOptions{}
}

// 应用选项
func (opts *getOptions) Apply(opt ...Option) {
	for _, o := range opt {
		o(opts)
	}
}

// 设置key
func WithGetKey(key []byte) Option {
	return func(opts *getOptions) {
		opts.Key = key
	}
}
