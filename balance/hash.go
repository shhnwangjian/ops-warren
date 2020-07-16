package balance

import (
	"fmt"
	"hash/crc32"
	"math/rand"
)

// 一致性hash算法
type HashBalance struct {
}

func init() {
	RegisterBalancer("hash", &HashBalance{})
}

func (p *HashBalance) DoBalance(instanceList []*Instance) (*Instance, error) {
	var defKey = fmt.Sprintf("%d", rand.Int())

	lens := len(instanceList)
	if lens == 0 {
		return nil, fmt.Errorf("no instance found")
	}

	crcTable := crc32.MakeTable(crc32.IEEE)
	hashVal := crc32.Checksum([]byte(defKey), crcTable)
	index := int(hashVal) % lens
	inst := instanceList[index]
	inst.CallNums++

	return inst, nil
}
