package bios

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
		}{Type: 0, Length: 26, Handle: 0},
		Formatted: []byte{1, 2, 0, 240, 3, 255, 144, 154, 233, 89, 0, 0, 31, 0, 3, 15, 2, 2, 255, 255, 32, 0},
		Strings:   []string{"Dell Inc.", "2.2.11", "06/13/2019"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
