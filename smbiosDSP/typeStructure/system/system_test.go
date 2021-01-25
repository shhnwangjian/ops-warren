package system

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
		}{Type: 1, Length: 27, Handle: 256},
		Formatted: []byte{1, 2, 0, 3, 68, 69, 76, 76, 82, 0, 16, 80, 128, 67, 180, 192, 79, 55, 90, 50, 6, 5, 4},
		Strings:   []string{"Dell Inc.", "PowerEdge R740", "4RPC7Z2", "PowerEdge", "SKU=NotProvided;ModelName=PowerEdge R740"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(fmt.Sprintf("%+v", res))
	t.Log(err)
}
