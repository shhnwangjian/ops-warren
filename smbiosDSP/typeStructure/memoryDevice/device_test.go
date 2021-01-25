package memoryDevice

import (
	"testing"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

func TestDeviceParseDell(t *testing.T) {
	ss := &dmi.SubSmBiosStructure{
		Header: struct {
			Type   uint8
			Length uint8
			Handle uint16
		}{Type: 17, Length: 84, Handle: 4352},
		Formatted: []byte{0, 16, 254, 255, 72, 0, 64, 0, 255, 127, 9, 1, 1, 0, 26, 128, 32, 106, 10, 2, 3, 4, 5, 2, 0,
			128, 0, 0, 96, 9, 176, 4, 176, 4, 176, 4, 3, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Strings: []string{"A1", "00AD063200AD", "2E9F863B", "011911A0", "HMA84GR7CJR4N-VK"},
	}
	t.Log(ss)
	res, err := Parse(ss)
	t.Log(res, err)
	t.Log(err)
}
