package cache

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

type CacheOperationalMode byte

const (
	CacheOperationalModeWriteThrough CacheOperationalMode = iota
	CacheOperationalModeWriteBack
	CacheOperationalModeVariesWithMemoryAddress
	CacheOperationalModeUnknown
)

func (c CacheOperationalMode) String() string {
	modes := [...]string{
		"Write Through",
		"Write Back",
		"Varies With Memory Address",
		"Unknown",
	}
	return modes[c]
}

type CacheLocation byte

const (
	CacheLocationInternal CacheLocation = iota
	CacheLocationExternal
	CacheLocationReserved
	CacheLocationUnknown
)

func (c CacheLocation) String() string {
	locations := [...]string{
		"Internal",
		"External",
		"Reserved",
		"Unknown",
	}
	return locations[c]
}

type CacheLevel byte

const (
	CacheLevel1 CacheLevel = iota
	CacheLevel2
	CacheLevel3
)

func (c CacheLevel) String() string {
	levels := [...]string{
		"Level1",
		"Level2",
		"Level3",
	}
	return levels[c]
}

type CacheConfiguration struct {
	Mode     CacheOperationalMode
	Enabled  bool
	Location CacheLocation
	Socketed bool
	Level    CacheLevel
}

func NewCacheConfiguration(u uint16) CacheConfiguration {
	var c CacheConfiguration
	c.Level = CacheLevel(byte(u & 0x7))
	c.Socketed = (u&0x10 == 1)
	c.Location = CacheLocation((u >> 5) & 0x3)
	fmt.Println(u & (0x7))
	c.Enabled = (u&(0x1<<7) == 1)
	c.Mode = CacheOperationalMode((u >> 8) & 0x7)
	return c
}

func (c CacheConfiguration) String() string {
	return fmt.Sprintf("\tLevel: %s\n"+
		"\tSocketed: %v\n"+
		"\tLocation: %s\n"+
		"\tEnabled: %v\n"+
		"\tOperational Mode: %s\t\t",
		c.Level,
		c.Socketed,
		c.Location,
		c.Enabled,
		c.Mode)
}

type CacheGranularity byte

const (
	CacheGranularity1K CacheGranularity = iota
	CacheGranularity64K
)

func (c CacheGranularity) String() string {
	grans := [...]string{
		"1K",
		"64K",
	}
	return grans[c]
}

type CacheSize struct {
	Granularity CacheGranularity
	Size        uint16
}

func NewCacheSize(u uint16) CacheSize {
	var c CacheSize
	c.Granularity = CacheGranularity(u >> 15)
	c.Size = u &^ (uint16(1) << 15)
	return c
}

func (c CacheSize) String() string {
	return fmt.Sprintf("%d * %s", c.Size, c.Granularity)
}

type CacheSRAMType uint16

const (
	CacheSRAMTypeOther CacheSRAMType = 1 << iota
	CacheSRAMTypeUnknown
	CacheSRAMTypeNonBurst
	CacheSRAMTypeBurst
	CacheSRAMTypePipelineBurst
	CacheSRAMTypeSynchronous
	CacheSRAMTypeAsynchronous
	CacheSRAMTypeReserved
)

func (c CacheSRAMType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Non-Burst",
		"Burst",
		"Pipeline Burst",
		"Synchronous",
		"Asynchronous",
		"Reserved",
	}
	return types[c/2]
}

type CacheSpeed byte

type CacheErrorCorrectionType byte

const (
	CacheErrorCorrectionTypeOther CacheErrorCorrectionType = 1 + iota
	CacheErrorCorrectionTypeUnknown
	CacheErrorCorrectionTypeNone
	CacheErrorCorrectionTypeParity
	CacheErrorCorrectionTypeSinglebitECC
	CacheErrorCorrectionTypeMultibitECC
)

func (c CacheErrorCorrectionType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC",
	}
	return types[c-1]
}

type CacheSystemCacheType byte

const (
	CacheSystemCacheTypeOther CacheSystemCacheType = 1 + iota
	CacheSystemCacheTypeUnknown
	CacheSystemCacheTypeInstruction
	CacheSystemCacheTypeData
	CacheSystemCacheTypeUnified
)

func (c CacheSystemCacheType) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Instruction",
		"Data",
		"Unified",
	}
	return types[c-1]
}

type CacheAssociativity byte

const (
	CacheAssociativityOther CacheAssociativity = 1 + iota
	CacheAssociativityUnknown
	CacheAssociativityDirectMapped
	CacheAssociativity2waySetAssociative
	CacheAssociativity4waySetAssociative
	CacheAssociativityFullyAssociative
	CacheAssociativity8waySetAssociative
	CacheAssociativity16waySetAssociative
	CacheAssociativity12waySetAssociative
	CacheAssociativity24waySetAssociative
	CacheAssociativity32waySetAssociative
	CacheAssociativity48waySetAssociative
	CacheAssociativity64waySetAssociative
	CacheAssociativity20waySetAssociative
)

func (c CacheAssociativity) String() string {
	caches := [...]string{
		"", // 0
		"Other",
		"Unknown",
		"Direct Mapped",
		"2-way Set-Associative",
		"4-way Set-Associative",
		"Fully Associative",
		"8-way Set-Associative",
		"16-way Set-Associative",
		"12-way Set-Associative",
		"24-way Set-Associative",
		"32-way Set-Associative",
		"48-way Set-Associative",
		"64-way Set-Associative",
		"20-way Set-Associative",
	}
	return caches[c]
}

type Information struct {
	dmi.Header
	SocketDesignation   string     `json:"socket_designation,omitempty"`
	Configuration       string     `json:"configuration,omitempty"`
	MaximumCacheSize    string     `json:"maximum_cache_size,omitempty"`
	InstalledSize       string     `json:"installed_size,omitempty"`
	SupportedSRAMType   string     `json:"supported_sram_type,omitempty"`
	CurrentSRAMType     string     `json:"current_sram_type,omitempty"`
	CacheSpeed          CacheSpeed `json:"cache_speed,omitempty"`
	ErrorCorrectionType string     `json:"error_correction_type,omitempty"`
	SystemCacheType     string     `json:"system_cache_type,omitempty"`
	Associativity       string     `json:"associativity,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"SocketDesignation":   b.SocketDesignation,
		"Configuration":       b.Configuration,
		"MaximumCacheSize":    b.MaximumCacheSize,
		"InstalledSize":       b.InstalledSize,
		"SupportedSRAMType":   b.SupportedSRAMType,
		"CurrentSRAMType":     b.CurrentSRAMType,
		"ErrorCorrectionType": b.ErrorCorrectionType,
		"SystemCacheType":     b.SystemCacheType,
		"Associativity":       b.Associativity,
	}
}

// ParseCache 缓存信息
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		Header:              s.Header,
		SocketDesignation:   s.GetString(0x0),
		CacheSpeed:          CacheSpeed(s.GetByte(0x0b)),
		ErrorCorrectionType: CacheErrorCorrectionType(s.GetByte(0xc)).String(),
		SystemCacheType:     CacheSystemCacheType(s.GetByte(0xd)).String(),
		Associativity:       CacheAssociativity(s.GetByte(0xe)).String(),
	}
	configuration, err := s.U16(0x01, 0x03)
	if err == nil {
		info.Configuration = NewCacheConfiguration(configuration).String()
	}
	maximumCacheSize, err := s.U16(0x03, 0x05)
	if err == nil {
		info.MaximumCacheSize = NewCacheSize(maximumCacheSize).String()
	}
	installedSize, err := s.U16(0x05, 0x07)
	if err == nil {
		info.InstalledSize = NewCacheSize(installedSize).String()
	}
	supportedSRAMType, err := s.U16(0x07, 0x09)
	if err == nil {
		info.SupportedSRAMType = CacheSRAMType(supportedSRAMType).String()
	}
	currentSRAMType, err := s.U16(0x09, 0x0b)
	if err == nil {
		info.CurrentSRAMType = CacheSRAMType(currentSRAMType).String()
	}
	return info.Map(), nil
}
