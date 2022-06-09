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
	HashKey []byte
	Target  string
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

// 指定目标
func WithTarget(name string) Option {
	return func(opts *getOptions) {
		opts.Target = name
	}
}

// 设置key
func WithHashKey(key string) Option {
	return func(opts *getOptions) {
		opts.HashKey = []byte(key)
	}
}
