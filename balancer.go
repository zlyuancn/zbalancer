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
	"fmt"
)

var NoInstanceErr = errors.New("no instance")

// 平衡器创建者
type BalancerCreator func() Balancer

// 平衡器
type Balancer interface {
	// 更新
	//
	// 如果实例发生变动, 调用此方法以告知平衡机
	Update(ins []interface{}, opt ...UpdateOption)
	// 获取一个实例, 如果实例总数为0会返回NoInstanceErr
	Get(opt ...Option) (interface{}, error)
}

// 平衡器类型
type BalancerType string

const (
	// 轮询
	RoundBalancer BalancerType = "round"
	// 加权随机
	WeightRandomBalancer BalancerType = "weight_random"
	// 加权hash
	WeightHashBalancer BalancerType = "weight_hash"
)

var balancerCreators = map[BalancerType]BalancerCreator{
	RoundBalancer:        newRoundBalancer,
	WeightRandomBalancer: newWeightRandomBalancer,
	WeightHashBalancer:   newWeightHashBalancer,
}

// 注册平衡器创建者, 应该在 NewBalancer 之前调用
func RegistryBalancerCreator(t BalancerType, creator BalancerCreator) {
	_, ok := balancerCreators[t]
	if ok {
		panic(fmt.Errorf("creator of BalancerType<%v> is registered", t))
	}
	balancerCreators[t] = creator
}

// 新实例化一个平衡器, 返回是否存在
func NewBalancer(t BalancerType) (Balancer, bool) {
	creator, ok := balancerCreators[t]
	if ok {
		return creator(), true
	}
	return nil, false
}
