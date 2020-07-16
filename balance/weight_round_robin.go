package balance

import (
	"fmt"
)

// 带权重的轮询调度算法
type WeightRoundRobinBalance struct {
	Index  int64
	Weight int64
}

func init() {
	RegisterBalancer("weight_roundrobin", &WeightRoundRobinBalance{})
}

func (p *WeightRoundRobinBalance) DoBalance(instanceList []*Instance) (*Instance, error) {
	lens := len(instanceList)
	if lens == 0 {
		return nil, fmt.Errorf("no instance found")
	}

	inst := p.GetInst(instanceList)
	inst.CallNums++

	return inst, nil
}

func (p *WeightRoundRobinBalance) GetInst(instanceList []*Instance) *Instance {
	gcd := getGCD(instanceList)
	for {
		p.Index = (p.Index + 1) % int64(len(instanceList))
		if p.Index == 0 {
			p.Weight = p.Weight - gcd
			if p.Weight <= 0 {
				p.Weight = getMaxWeight(instanceList)
				if p.Weight == 0 {
					return &Instance{}
				}
			}
		}

		if instanceList[p.Index].Weight >= p.Weight {
			return instanceList[p.Index]
		}
	}
}

// gcd 计算两个数的最大公约数
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// getGCD 计算多个数的最大公约数
func getGCD(instanceList []*Instance) int64 {
	var weights []int64

	for _, instance := range instanceList {
		weights = append(weights, instance.Weight)
	}

	g := weights[0]
	for i := 1; i < len(weights)-1; i++ {
		oldGcd := g
		g = gcd(oldGcd, weights[i])
	}
	return g
}

// getMaxWeight 获取最大权重
func getMaxWeight(instanceList []*Instance) int64 {
	var max int64 = 0
	for _, instance := range instanceList {
		if instance.Weight >= int64(max) {
			max = instance.Weight
		}
	}

	return max
}
