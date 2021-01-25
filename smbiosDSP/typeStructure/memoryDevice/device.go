package memoryDevice

import (
	"strings"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

type MemoryDeviceFormFactor byte

const (
	MemoryDeviceFormFactorOther MemoryDeviceFormFactor = 1 + iota
	MemoryDeviceFormFactorUnknown
	MemoryDeviceFormFactorSIMM
	MemoryDeviceFormFactorSIP
	MemoryDeviceFormFactorChip
	MemoryDeviceFormFactorDIP
	MemoryDeviceFormFactorZIP
	MemoryDeviceFormFactorProprietaryCard
	MemoryDeviceFormFactorDIMM
	MemoryDeviceFormFactorTSOP
	MemoryDeviceFormFactorRowofchips
	MemoryDeviceFormFactorRIMM
	MemoryDeviceFormFactorSODIMM
	MemoryDeviceFormFactorSRIMM
	MemoryDeviceFormFactorFB_DIMM
)

func (m MemoryDeviceFormFactor) String() string {
	factors := [...]string{
		"Other",
		"Unknown",
		"SIMM",
		"SIP",
		"Chip",
		"DIP",
		"ZIP",
		"Proprietary Card",
		"DIMM",
		"TSOP",
		"Row of chips",
		"RIMM",
		"SODIMM",
		"SRIMM",
		"FB-DIMM",
	}
	return factors[m-1]
}

type MemoryDeviceType byte

const (
	MemoryDeviceTypeOther MemoryDeviceType = 1 + iota
	MemoryDeviceTypeUnknown
	MemoryDeviceTypeDRAM
	MemoryDeviceTypeEDRAM
	MemoryDeviceTypeVRAM
	MemoryDeviceTypeSRAM
	MemoryDeviceTypeRAM
	MemoryDeviceTypeROM
	MemoryDeviceTypeFLASH
	MemoryDeviceTypeEEPROM
	MemoryDeviceTypeFEPROM
	MemoryDeviceTypeEPROM
	MemoryDeviceTypeCDRAM
	MemoryDeviceType3DRAM
	MemoryDeviceTypeSDRAM
	MemoryDeviceTypeSGRAM
	MemoryDeviceTypeRDRAM
	MemoryDeviceTypeDDR
	MemoryDeviceTypeDDR2
	MemoryDeviceTypeDDR2FB_DIMM
	_
	_
	_
	MemoryDeviceTypeDDR3
	MemoryDeviceTypeFBD2
	MemoryDeviceTypeDDR4
	MemoryDeviceTypeLPDDR
	MemoryDeviceTypeLPDDR2
	MemoryDeviceTypeLPDDR3
	MemoryDeviceTypeLPDDR4
)

func (m MemoryDeviceType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"DRAM",
		"EDRAM",
		"VRAM",
		"SRAM",
		"RAM",
		"ROM",
		"FLASH",
		"EEPROM",
		"FEPROM",
		"EPROM",
		"CDRAM",
		"3DRAM",
		"SDRAM",
		"SGRAM",
		"RDRAM",
		"DDR",
		"DDR2",
		"DDR2 FB-DIMM",
		"Reserved1",
		"Reserved2",
		"Reserved3",
		"DDR3",
		"FBD2",
		"DDR4",
		"LPDDR",
		"LPDDR2",
		"LPDDR3",
		"LPDDR4",
	}
	return types[m-1]
}

type MemoryDeviceTypeDetail uint16

const (
	MemoryDeviceTypeDetailReserved MemoryDeviceTypeDetail = 1 + iota
	MemoryDeviceTypeDetailOther
	MemoryDeviceTypeDetailUnknown
	MemoryDeviceTypeDetailFast_paged
	MemoryDeviceTypeDetailStaticcolumn
	MemoryDeviceTypeDetailPseudo_static
	MemoryDeviceTypeDetailRAMBUS
	MemoryDeviceTypeDetailSynchronous
	MemoryDeviceTypeDetailCMOS
	MemoryDeviceTypeDetailEDO
	MemoryDeviceTypeDetailWindowDRAM
	MemoryDeviceTypeDetailCacheDRAM
	MemoryDeviceTypeDetailNon_volatile
	MemoryDeviceTypeDetailRegisteredBuffered
	MemoryDeviceTypeDetailUnbufferedUnregistered
	MemoryDeviceTypeDetailLRDIMM
)

func (m MemoryDeviceTypeDetail) String() string {
	details := [...]string{
		"Reserved",
		"Other",
		"Unknown",
		"Fast-paged",
		"Static column",
		"Pseudo-static",
		"RAMBUS",
		"Synchronous",
		"CMOS",
		"EDO",
		"Window DRAM",
		"Cache DRAM",
		"Non-volatile",
		"Registered (Buffered)",
		"Unbuffered (Unregistered)",
		"LRDIMM",
	}

	d := []string{}
	for i := 1; i <= 16; i++ {
		if m&(1<<uint(i)) != 0 {
			d = append(d, details[i])
		}
	}

	return strings.Join(d, ",")
}

// MemoryDevice 内存设备
type Information struct {
	dmi.Header
	PhysicalMemoryArrayHandle  uint16
	ErrorInformationHandle     uint16
	TotalWidth                 uint16
	DataWidth                  uint16
	Size                       uint16
	FormFactor                 string
	DeviceSet                  byte
	DeviceLocator              string
	BankLocator                string
	Type                       string
	TypeDetail                 string
	Speed                      uint16
	Manufacturer               string
	SerialNumber               string
	AssetTag                   string
	PartNumber                 string
	Attributes                 byte
	ExtendedSize               uint32
	ConfiguredMemoryClockSpeed uint16
	MinimumVoltage             uint16
	MaximumVoltage             uint16
	ConfiguredVoltage          uint16
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Manufacturer":  b.Manufacturer,
		"FormFactor":    b.FormFactor,
		"DeviceLocator": b.DeviceLocator,
		"SerialNumber":  b.SerialNumber,
		"AssetTag":      b.AssetTag,
		"BankLocator":   b.BankLocator,
		"Type":          b.Type,
		"TypeDetail":    b.TypeDetail,
		"PartNumber":    b.PartNumber,
	}
}

// Parse Memory Device
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		FormFactor:    MemoryDeviceFormFactor(s.GetByte(0x0a)).String(),
		DeviceSet:     s.GetByte(0x0b),
		DeviceLocator: s.GetString(0xc),
		BankLocator:   s.GetString(0xd),
		Type:          MemoryDeviceType(s.GetByte(0xe)).String(),
		Manufacturer:  s.GetString(0x13),
		SerialNumber:  s.GetString(0x14),
		PartNumber:    s.GetString(0x16),
		Attributes:    s.GetByte(0x17),
	}
	physicalMemoryArrayHandle, err := s.U16(0x00, 0x02)
	if err == nil {
		info.PhysicalMemoryArrayHandle = physicalMemoryArrayHandle
	}
	errorInformationHandle, err := s.U16(0x02, 0x04)
	if err == nil {
		info.ErrorInformationHandle = errorInformationHandle
	}
	totalWidth, err := s.U16(0x04, 0x06)
	if err == nil {
		info.TotalWidth = totalWidth
	}
	dataWidth, err := s.U16(0x06, 0x08)
	if err == nil {
		info.DataWidth = dataWidth
	}
	size, err := s.U16(0x08, 0x0a)
	if err == nil {
		info.Size = size
	}
	typeDetail, err := s.U16(0xf, 0x11)
	if err == nil {
		info.TypeDetail = MemoryDeviceTypeDetail(typeDetail).String()
	}
	speed, err := s.U16(0x11, 0x13)
	if err == nil {
		info.Speed = speed
	}
	extendedSize, err := s.U32(0x18, 0x1c)
	if err == nil {
		info.ExtendedSize = extendedSize
	}
	configuredMemoryClockSpeed, err := s.U16(0x1c, 0x1e)
	if err == nil {
		info.ConfiguredMemoryClockSpeed = configuredMemoryClockSpeed
	}
	minimumVoltage, err := s.U16(0x1e, 0x20)
	if err == nil {
		info.MinimumVoltage = minimumVoltage
	}
	maximumVoltage, err := s.U16(0x20, 0x22)
	if err == nil {
		info.MaximumVoltage = maximumVoltage
	}
	configuredVoltage, err := s.U16(0x22, 0x24)
	if err == nil {
		info.ConfiguredVoltage = configuredVoltage
	}
	return info.Map(), nil
}
