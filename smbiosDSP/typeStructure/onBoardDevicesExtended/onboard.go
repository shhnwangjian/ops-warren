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

// ExtendedInformation 扩展设备信息 (smbios >= 2.6)
// 参考文档: https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf
// 7.42 Onboard Devices Extended Information (Type 41)
type ExtendedInformation struct {
	dmi.Header
	ReferenceDesignation string `json:"reference_designation,omitempty"`
	DeviceType           string `json:"device_type,omitempty"`
	DeviceStatus         string `json:"device_status,omitempty"`
	DeviceTypeInstance   byte   `json:"device_type_instance,omitempty"`
	SegmentGroupNumber   uint16 `json:"segment_group_number,omitempty"`
	BusNumber            byte   `json:"bus_number,omitempty"`
	DeviceFunctionNumber byte   `json:"device_function_number,omitempty"`
}

// SlotSegment 主板设备对应的slot
func (o ExtendedInformation) SlotSegment() string {
	if o.SegmentGroupNumber == 0xFFFF || o.BusNumber == 0xFF || o.DeviceFunctionNumber == 0xFF {
		return "Not of types PCI/AGP/PCI-X/PCI-Express"
	}
	return fmt.Sprintf("Bus Address: %04x:%02x:%02x.%x",
		o.SegmentGroupNumber,
		o.BusNumber,
		o.DeviceFunctionNumber>>3,
		o.DeviceFunctionNumber&0x7)
}

// ParseType41 解析
func Parse(s *dmi.SubSmBiosStructure) (info *ExtendedInformation, err error) {

	if s.Type() != 41 {
		return nil, fmt.Errorf("ParseType41 only parse type 41 data, but now: %d", s.Type())
	}

	info = &ExtendedInformation{
		Header:               s.Header,
		ReferenceDesignation: s.GetString(0x0),
		DeviceType:           ExtendedInformationType(s.GetByte(0x01) & 127).String(),
		DeviceStatus:         DeviceStatus(s.GetByte(0x01) >> 7).String(),
		DeviceTypeInstance:   s.GetByte(0x02),
		BusNumber:            s.GetByte(0x05),
		DeviceFunctionNumber: s.GetByte(0x06),
	}
	segmentGroupNumber, err := s.U16(0x03, 0x05)
	if err == nil {
		info.SegmentGroupNumber = segmentGroupNumber
	}
	return info, nil
}
