
# 平衡器

- [x] 轮询
- [x] 加权随机
- [ ] 加权hash
- [ ] 加权一致性hash环

# 示例

```go
b, _ := zbalancer.NewBalancer(zbalancer.RoundBalancer) // 创建一个轮询平衡器
b.Update([]interface{}{"nodeA", "nodeB", "nodeC"}) // 重设节点

node := b.Get() // 获取节点
```
