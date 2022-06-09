package zbalancer

// 目标选择器
type targetSelector struct {
	insMap map[string]Instance // 用于直接获取目标
}

func newTargetSelector() *targetSelector {
	return new(targetSelector)
}

func (t *targetSelector) Update(instances []Instance) {
	t.insMap = make(map[string]Instance, len(instances))
	for _, ins := range instances {
		t.insMap[ins.Name()] = ins
	}
}

func (t *targetSelector) Get(target string) (Instance, error) {
	ins, ok := t.insMap[target]
	if !ok {
		return nil, NoInstanceErr
	}
	return ins, nil
}
