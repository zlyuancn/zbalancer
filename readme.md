
# 平衡器

- [x] [RoundBalancer](./round.go) 轮询
  > 按顺序获取实例.
- [x] [WeightRandomBalancer](./weight_random.go) 加权随机
  > 每个实例有不同权重, 获取时随机选择一个实例, 权重越高被选取的机会越大.
- [x] 加权hash
  > 每个实例有不同的权重, 获取时根据提供的key计算hash值然后对总权重求余, 余数计算所在实例, 权重越高被选取的机会越大.
- [ ] 加权一致性hash环
  > 每个实例有不同的权重, 权重值可以理解为每个实例的分片数, 每个分片计算hash值落在一个环上. 获取时根据提供的key计算hash值然后得出落在环的一个点上, 由这个点得出是哪个实例的分片进而知道是哪个实例.

# 示例

## 轮询

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.RoundBalancer) // 创建一个轮询平衡器
balancer.Update([]interface{}{"nodeA", "nodeB", "nodeC"}) // 重设节点

node, _ := balancer.Get() // 获取节点
fmt.Println(node)
```

## 加权随机

```go
balancer, _ := zbalancer.NewBalancer(zbalancer.WeightRandomBalancer) // 创建一个权重随机平衡器
balancer.Update(
    []interface{}{"nodeA", "nodeB", "nodeC"},      // 重设节点
    zbalancer.WithUpdateWeights([]uint8{1, 2, 3}), // 设置权重
)

node, _ := balancer.Get() // 获取节点
fmt.Println(node)
```
