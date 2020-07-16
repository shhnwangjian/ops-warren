package balance

import (
	"fmt"
)

var (
	balanceMgr = BalancerManager{
		allBalance: make(map[string]Balancer),
	}
)

// 负载均衡
type Balancer interface {
	DoBalance([]*Instance) (*Instance, error)
}

// 负载均衡管理器
type BalancerManager struct {
	allBalance map[string]Balancer
}

func (p *BalancerManager) register(balanceType string, b Balancer) {
	p.allBalance[balanceType] = b
}

func RegisterBalancer(balanceType string, b Balancer) {
	balanceMgr.register(balanceType, b)
}

func DoBalance(balanceType string, instanceList []*Instance) (*Instance, error) {
	Balancer, ok := balanceMgr.allBalance[balanceType]
	if !ok {
		return nil, fmt.Errorf("not found %s balancer", balanceType)
	}
	return Balancer.DoBalance(instanceList)
}
