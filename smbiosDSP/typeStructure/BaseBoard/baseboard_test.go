package baseboard

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
		}{Type: 2, Length: 8, Handle: 512},
		Formatted: []byte{1, 2, 3, 4},
		Strings:   []string{"Dell Inc.", "01YM03", "A02", ".4RPC7Z2.CNIVC009740870."},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
