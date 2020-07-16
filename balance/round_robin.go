package balance

import (
	"fmt"
)

// 轮询调度算法
type RoundRobinBalance struct {
	curIndex int
}

func init() {
	RegisterBalancer("roundrobin", &RoundRobinBalance{})
}

func (p *RoundRobinBalance) DoBalance(instanceList []*Instance) (*Instance, error) {
	lens := len(instanceList)
	if lens == 0 {
		return nil, fmt.Errorf("no instance found")
	}

	if p.curIndex >= lens {
		p.curIndex = 0
	}
	inst := instanceList[p.curIndex]

	p.curIndex = (p.curIndex + 1) % lens

	inst.CallNums++

	return inst, nil
}
