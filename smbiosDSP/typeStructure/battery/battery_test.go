package battery

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
		}{Type: 22, Length: 26, Handle: 44},
		Formatted: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00},
		Strings:   []string{"Fake", "-Virtual Battery 0-", "08/08/2010", "Battery 0", "CRB Battery 0", "LithiumPolymer"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
}
