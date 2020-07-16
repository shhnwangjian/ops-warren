package balance

import (
	"errors"
	"math/rand"
)

// 随机算法
type RandomBalance struct {
}

func init() {
	RegisterBalancer("random", &RandomBalance{})
}

func (p *RandomBalance) DoBalance(instanceList []*Instance) (*Instance, error) {
	lens := len(instanceList)
	if lens == 0 {
		return nil, errors.New("no instance found")
	}

	index := rand.Intn(lens)
	inst := instanceList[index]
	inst.CallNums++

	return inst, nil
}
