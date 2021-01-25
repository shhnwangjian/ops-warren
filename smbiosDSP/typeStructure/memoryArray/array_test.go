package memoryArray

import (
	"testing"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

func TestParseDell(t *testing.T) {
	ss := &dmi.SubSmBiosStructure{
		Header: struct {
			Type   uint8
			Length uint8
			Handle uint16
		}{Type: 16, Length: 23, Handle: 4096},
		Formatted: []byte{3, 3, 6, 0, 0, 0, 128, 254, 255, 24, 0, 0, 0, 0, 0, 128, 7, 0, 0},
		Strings:   []string{},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
