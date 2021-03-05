/*
-------------------------------------------------
   Author :       zlyuancn
   date：         2021/3/5
   Description :
-------------------------------------------------
*/

package zbalancer

// 平衡器
type Balancer interface {
	// 更新
	//
	// 如果实例发生变动, 调用此方法以告知平衡机
	Update(ins []interface{}, opt ...Option)
	// 获取一个实例, 如果实例总数为0会返回nil
	Get() interface{}
}

// 平衡器类型
type BalancerType string

const (
	// 轮询
	RoundBalancer BalancerType = "round"
	// 加权随机
	WeightRandomBalancer BalancerType = "weight_random"
)

var balancers = map[BalancerType]Balancer{
	RoundBalancer:        newRoundBalancer(),
	WeightRandomBalancer: newWeightRandomBalancer(),
}

// 注册平衡器, 应该在 NewBalancer 之前调用
func RegistryBalancer(t BalancerType, balancer Balancer) {
	balancers[t] = balancer
}

// 新实例化一个平衡器, 返回是否存在
func NewBalancer(t BalancerType) (Balancer, bool) {
	b, ok := balancers[t]
	return b, ok
}
