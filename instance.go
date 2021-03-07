/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/7
   Description :
-------------------------------------------------
*/

package zbalancer

type Instance interface {
	// 获取实例
	Instance() interface{}
	// 获取实例名, 部分平衡器要求要有name, 否则会产生异常
	Name() string
	// 获取权重值, 部分平衡器根据权重值获取实例, 权重值为0的实例会被忽略
	Weight() uint8
}

type instanceCli struct {
	instance interface{}
	name     string
	weight   uint8
}

func (i *instanceCli) Instance() interface{} {
	return i.instance
}

func (i *instanceCli) Name() string {
	return i.name
}

func (i *instanceCli) Weight() uint8 {
	return i.weight
}

// 设置实例名
func (i *instanceCli) SetName(name string) *instanceCli {
	i.name = name
	return i
}

// 设置权重
func (i *instanceCli) SetWeight(weight uint8) *instanceCli {
	i.weight = weight
	return i
}

func NewInstance(instance interface{}) *instanceCli {
	return &instanceCli{
		instance: instance,
		weight:   1,
	}
}
