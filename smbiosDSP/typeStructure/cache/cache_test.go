package cache

import (
	"testing"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

func TestParseProcessorDell(t *testing.T) {
	ss := &dmi.SubSmBiosStructure{
		Header: struct {
			Type   uint8
			Length uint8
			Handle uint16
		}{Type: 7, Length: 27, Handle: 1792},
		Formatted: []byte{0, 128, 1, 0, 3, 0, 3, 2, 0, 2, 0, 0, 4, 5, 7, 0, 0, 0, 0, 0, 0, 0, 0},
		Strings:   []string{},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
