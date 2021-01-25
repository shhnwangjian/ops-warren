package bios

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

// Information bios信息
type Information struct {
	dmi.Header
	Vendor                 string `json:"vendor,omitempty"`
	BIOSVersion            string `json:"bios_version,omitempty"`
	StartingAddressSegment string `json:"starting_address_segment,omitempty"`
	ReleaseDate            string `json:"release_date,omitempty"`
	RomSize                string `json:"rom_size,omitempty"`
	RuntimeSize            string `json:"runtime_size,omitempty"`
	Characteristics        string `json:"characteristics,omitempty"`
	CharacteristicsExt1    string `json:"characteristics_ext_1,omitempty"`
	CharacteristicsExt2    string `json:"characteristics_ext_2,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Vendor":                 b.Vendor,
		"BIOSVersion":            b.BIOSVersion,
		"StartingAddressSegment": b.StartingAddressSegment,
		"ReleaseDate":            b.ReleaseDate,
		"RomSize":                b.RomSize,
		"RuntimeSize":            b.RuntimeSize,
		"Characteristics":        b.Characteristics,
		"CharacteristicsExt1":    b.CharacteristicsExt1,
		"CharacteristicsExt2":    b.CharacteristicsExt2,
	}
}

func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	if len(s.Strings) == 0 {
		return make(map[string]string, 0), fmt.Errorf("no data")
	}
	info := &Information{
		Header:              s.Header,
		Vendor:              s.GetString(0x0),
		BIOSVersion:         s.GetString(0x1),
		ReleaseDate:         s.GetString(0x4),
		RomSize:             RomSize(64 * (s.GetByte(0x05) + 1)).String(),
		CharacteristicsExt1: Ext1(s.GetByte(0x08)).String(),
		CharacteristicsExt2: Ext2(s.GetByte(0x09)).String(),
	}
	address, err := s.U16(0x02, 0x04)
	if err == nil {
		info.StartingAddressSegment = fmt.Sprintf("0x%4X0", address)
		info.RuntimeSize = RuntimeSize((uint(0x10000) - uint(address)) << 4).String()
	}
	characteristic, err := s.U64(0x06, 0x08)
	if err == nil {
		info.Characteristics = Characteristics(characteristic).String()
	}
	return info.Map(), nil
}
