package onboardDevice

import (
	"fmt"
	"testing"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

func TestParseDell(t *testing.T) {
	ss := &dmi.SubSmBiosStructure{
		Header: struct {
			Type   uint8
			Length uint8
			Handle uint16
		}{Type: 41, Length: 11, Handle: 10496},
		Formatted: []byte{1, 133, 1, 0, 0, 24},
		Strings:   []string{"Integrated NIC 1"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
