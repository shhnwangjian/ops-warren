package system

import (
	"bytes"
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

// WakeUpType 唤醒类型
type WakeUpType byte

const (
	Reserved WakeUpType = iota
	Other
	Unknown
	APMTimer
	ModemRing
	LANRemote
	PoAerSwitch
	PCIPME
	ACPowerRestored
)

func (w WakeUpType) String() string {
	types := [...]string{
		"Reserved", /* 0x00 */
		"Other",
		"Unknown",
		"APM Timer",
		"Modem Ring",
		"LAN Remote",
		"Power Switch",
		"PCI PME#",
		"AC Power Restored", /* 0x08 */
	}
	return types[w]
}

// Information 系统信息
type Information struct {
	dmi.Header
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	Version      string `json:"version,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	UUID         string `json:"uuid,omitempty"`
	WakeUpType   string `json:"wake_up_type,omitempty"`
	SKUNumber    string `json:"sku_number,omitempty"`
	Family       string `json:"family,omitempty"`
}

func uuid(data []byte, ver string) string {
	if bytes.Index(data, []byte{0x00}) != -1 {
		return "Not present"
	}

	if bytes.Index(data, []byte{0xFF}) != -1 {
		return "Not settable"
	}

	if ver > "2.6" {
		return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
			data[3], data[2], data[1], data[0], data[5], data[4], data[7], data[6],
			data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15])
	}
	return fmt.Sprintf("%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7],
		data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15])
}

// Parse 解析smbios struct数据
func Parse(s *dmi.SubSmBiosStructure) (info *Information, err error) {

	info = &Information{
		Header:       s.Header,
		Manufacturer: s.GetString(0x0),
		ProductName:  s.GetString(0x1),
		Version:      s.GetString(0x2),
		SerialNumber: s.GetString(0x3),
		UUID:         uuid(s.GetBytes(0x04, 0x14), s.GetString(2)),
		WakeUpType:   WakeUpType(s.GetByte(0x14)).String(),
		SKUNumber:    s.GetString(0x15),
		Family:       s.GetString(0x16),
	}

	return info, nil
}
