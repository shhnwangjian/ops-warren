package onboardDevice

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

// ExtendedInformationType 主板上的设备类型
type ExtendedInformationType byte

const (
	ExtendedInformationTypeOther ExtendedInformationType = 1 + iota
	ExtendedInformationTypeUnknown
	ExtendedInformationTypeVideo
	ExtendedInformationTypeSCSIController
	ExtendedInformationTypeEthernet
	ExtendedInformationTypeTokenRing
	ExtendedInformationTypeSound
	ExtendedInformationTypePATAController
	ExtendedInformationTypeSATAController
	ExtendedInformationTypeSASController
)

func (o ExtendedInformationType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Video",
		"SCSI Controller",
		"Ethernet",
		"Token Ring",
		"Sound",
		"PATA Controller",
		"SATA Controller",
		"SAS Controller",
	}
	return types[o-1]
}

// DeviceStatus 设备状态
type DeviceStatus byte

const (
	DeviceDisabled DeviceStatus = 0
	DeviceEnabled               = 1
)

func (d DeviceStatus) String() string {
	types := [...]string{
		"Disabled",
		"Enabled",
	}
	return types[d]
}

// Information 设备信息 (smbios < 2.6)
// 参考: https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf
//      7.11 On Board Devices Information (Type 10, Obsolete)
type Information struct {
	dmi.Header
	Devices []Device `json:"devices,omitempty"`
}

// Device 设备
type Device struct {
	Description  string `json:"description,omitempty"`
	DeviceStatus string `json:"device_status,omitempty"`
	DeviceType   string `json:"device_type,omitempty"`
}

// ParseType10 解析
func Parse(s *dmi.SubSmBiosStructure) (info *Information, err error) {

	if s.Type() != 10 {
		return nil, fmt.Errorf("ParseType10 only parse type 10 data, but now: %d", s.Type())
	}

	data := s.Formatted

	var devices []Device
	for i := 0x0; i <= len(data)-1; i += 2 {
		d := Device{
			DeviceStatus: DeviceStatus(data[i] >> 7).String(),
			DeviceType:   ExtendedInformationType(data[i] & 127).String(),
			Description:  s.GetString(i + 1),
		}
		devices = append(devices, d)
	}

	info = &Information{
		Header:  s.Header,
		Devices: devices,
	}

	return info, nil
}
