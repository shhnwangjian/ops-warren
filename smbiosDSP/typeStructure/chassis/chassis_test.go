package chassis

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
		}{Type: 3, Length: 22, Handle: 768},
		Formatted: []byte{1, 151, 0, 2, 0, 3, 3, 3, 2, 0, 0, 0, 0, 2, 0, 0, 0, 0},
		Strings:   []string{"Dell Inc.", "4RPC7Z2"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
