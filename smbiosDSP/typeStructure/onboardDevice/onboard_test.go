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
		}{Type: 10, Length: 11, Handle: 10498},
		Formatted: []byte{0x85, 0x1, 0x85, 0x2, 0x85, 0x3},
		Strings:   []string{"d1", "d2", "d3"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
