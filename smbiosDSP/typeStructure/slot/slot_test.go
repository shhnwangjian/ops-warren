package slot

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
		}{Type: 9, Length: 17, Handle: 2304},
		Formatted: []byte{1, 182, 11, 4, 4, 1, 0, 4, 1, 0, 0, 59, 0},
		Strings:   []string{"PCIe Slot 1"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
