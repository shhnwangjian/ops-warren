package port

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
		}{Type: 8, Length: 9, Handle: 2048},
		Formatted: []byte{1, 18, 0, 0, 16},
		Strings:   []string{"Internal USB port 1"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
