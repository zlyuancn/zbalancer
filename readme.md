# 简单易懂的现代平衡器

- [RoundBalancer](./round.go) 轮询
  > 按顺序获取实例.
- [WeightRandomBalancer](./weight_random.go) 加权随机
  > 每个实例有不同权重, 获取时随机选择一个实例, 权重越高被选取的机会越大.
- [WeightHashBalancer](./weight_hash.go) 加权hash
  > 每个实例有不同的权重, 获取时根据提供的key计算hash值然后对总权重求余, 余数计算所在实例, 权重越高被选取的机会越大.
- [WeightConsistentHashBalancer](./weight_consistent_hash.go)加权一致性hash环
  > 每个实例有不同的权重, 权重值可以理解为每个实例的分片数, 每个分片计算hash值落在一个环上. 获取时根据提供的key计算hash值然后得出落在环的一个点上, 由这个点得出是哪个实例的分片进而知道是哪个实例.

# 示例

## 轮询

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.RoundBalancer) // 创建一个轮询平衡器
balancer.Update( // 重设节点
    zbalancer.NewInstance("nodeA"),
    zbalancer.NewInstance("nodeB"),
    zbalancer.NewInstance("nodeC"),
)

node, _ := balancer.Get() // 获取节点
fmt.Println(node.Instance())
```

## 加权随机

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.WeightRandomBalancer) // 创建一个加权随机平衡器
balancer.Update( // 重设节点
    zbalancer.NewInstance("nodeA").SetWeight(1), // 设置权重
    zbalancer.NewInstance("nodeB").SetWeight(2),
    zbalancer.NewInstance("nodeC").SetWeight(3),
)

node, _ := balancer.Get() // 获取节点
fmt.Println(node.Instance())
```

## 加权hash

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.WeightHashBalancer) // 创建一个加权hash平衡器
balancer.Update(
    zbalancer.NewInstance("nodeA").SetWeight(1), // 设置权重
    zbalancer.NewInstance("nodeB").SetWeight(2),
    zbalancer.NewInstance("nodeC").SetWeight(3),
)

node, _ := balancer.Get(zbalancer.WithKey([]byte("hello"))) // 根据key获取节点
fmt.Println(node.Instance())
```

## 加权一致性hash环

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.WeightConsistentHashBalancer) // 创建一个加权一致性hash平衡器
balancer.Update(
    // 设置实例名和权重. 注意: 必须设置实例名, 且每个实例的实例名不能相同, 否则会导致异常
    zbalancer.NewInstance("nodeA").SetName("nodeA").SetWeight(1),
    zbalancer.NewInstance("nodeB").SetName("nodeB").SetWeight(2),
    zbalancer.NewInstance("nodeC").SetName("nodeC").SetWeight(3),
)

node, _ := balancer.Get(zbalancer.WithKey([]byte("hello"))) // 根据key获取节点
fmt.Println(node.Instance())
```
