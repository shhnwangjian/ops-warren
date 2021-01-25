package processor

import (
	"fmt"
	"testing"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

func TestParseProcessorDell(t *testing.T) {
	ss := &dmi.SubSmBiosStructure{
		Header: struct {
			Type   uint8
			Length uint8
			Handle uint16
		}{Type: 4, Length: 48, Handle: 1024},
		Formatted: []byte{1, 3, 179, 2, 84, 6, 5, 0, 255, 251, 235, 191, 3, 146, 160, 40, 160, 15, 252, 8, 65, 38, 0,
			7, 1, 7, 2, 7, 0, 0, 0, 12, 12, 24, 252, 0, 179, 0, 12, 0, 12, 0, 24, 0},
		Strings: []string{"CPU1", "Intel", "Intel(R) Xeon(R) Gold 5118 CPU @ 2.30GHz"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
