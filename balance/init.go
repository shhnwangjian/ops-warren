package balance

import (
	"fmt"
)

type Instance struct {
	Ip       string
	Port     int64
	Weight   int64
	CallNums int64
}

func NewInstance(ip string, port int64, w int64) *Instance {
	return &Instance{
		Ip:       ip,
		Port:     port,
		Weight:   w, // 权重
		CallNums: 0, // 调用次数
	}
}

func (i *Instance) GetPort() int64 {
	return i.Port
}

func (i *Instance) GetHost() string {
	return i.Ip
}

func (i *Instance) GetAddr() string {
	return fmt.Sprintf("%s:%d", i.Ip, i.Port)
}

func (i *Instance) GetCallTimes() int64 {
	return i.CallNums
}

func (i *Instance) GetResult() string {
	return fmt.Sprintf("%s:%d;call nums:%d", i.Ip, i.Port, i.CallNums)
}
