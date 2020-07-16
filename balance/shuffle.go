package balance

import (
	"errors"
	"math/rand"
	"time"
)

// 洗牌算法
type Shuffle2Balance struct {
}

func init() {
	RegisterBalancer("shuffle", &Shuffle2Balance{})
}

func (p *Shuffle2Balance) DoBalance(instanceList []*Instance) (*Instance, error) {
	lens := len(instanceList)
	if lens == 0 {
		return nil, errors.New("no instance found")
	}

	rand.Seed(time.Now().UnixNano())

	for i := lens; i > 0; i-- {
		lastIdx := i - 1
		idx := rand.Intn(i)
		instanceList[lastIdx], instanceList[idx] = instanceList[idx], instanceList[lastIdx]
	}

	inst := instanceList[0]
	inst.CallNums++

	return inst, nil
}
