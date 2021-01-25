package baseboard

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

// FeatureFlags 主板功能标签
type FeatureFlags byte

// Baseboard feature flags
const (
	HostingBoard FeatureFlags = 1 << iota
	AtLeastOneDaughter
	Removable
	Repleaceable
	HotSwappable
	//FeatureFlagsReserved = 000b
)

func (f FeatureFlags) String() string {
	fmt.Printf("xxx,%d", f)
	features := [...]string{
		"Board is a hosting board", /* 0 */
		"Board requires at least one daughter board",
		"Board is removable",
		"Board is replaceable",
		"Board is hot swappable", /* 4 */
	}
	var s string
	for i := uint32(0); i < 5; i++ {
		if f&(1<<i) != 0 {
			s += "\n\t\t" + features[i]
		}
	}
	return s
}

// Type 主板类型
type Type byte

const (
	Unknown Type = 1 + iota
	Other
	ServerBlade
	ConnectivitySwitch
	ManagementModule
	ProcessorModule
	IOModule
	MemModule
	DaughterBoard
	Motherboard
	ProcessorMemmoryModule
	ProcessorIOModule
	InterconnectBoard
)

func (b Type) String() string {
	types := [...]string{
		"Unknown", /* 0x01 */
		"Other",
		"Server Blade",
		"Connectivity Switch",
		"System Management Module",
		"Processor Module",
		"I/O Module",
		"Memory Module",
		"Daughter Board",
		"Motherboard",
		"Processor+Memory Module",
		"Processor+I/O Module",
		"Interconnect Board", /* 0x0D */
	}
	if b > Unknown && b < InterconnectBoard {
		return types[b-1]
	}
	return "Out Of Spec"
}

// Information 主板信息
type Information struct {
	dmi.Header
	Manufacturer      string `json:"manufacturer,omitempty"`
	ProductName       string `json:"product_name,omitempty"`
	Version           string `json:"version,omitempty"`
	SerialNumber      string `json:"serial_number,omitempty"`
	AssetTag          string `json:"asset_tag,omitempty"`
	FeatureFlags      string `json:"feature_flags,omitempty"`
	LocationInChassis string `json:"location_in_chassis,omitempty"`
	//ChassisHandle                  uint16 `json:"chassis_handle,omitempty"`
	BoardType string `json:"board_type,omitempty"`
	//NumberOfContainedObjectHandles byte   `json:"number_of_contained_object_handles,omitempty"`
	//ContainedObjectHandles         []byte `json:"contained_object_handles,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Manufacturer":      b.Manufacturer,
		"ProductName":       b.ProductName,
		"Version":           b.Version,
		"SerialNumber":      b.SerialNumber,
		"AssetTag":          b.AssetTag,
		"FeatureFlags":      b.FeatureFlags,
		"LocationInChassis": b.LocationInChassis,
		"BoardType":         b.BoardType,
	}
}

// Parse 解析
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		Header:            s.Header,
		Manufacturer:      s.GetString(0x0),
		ProductName:       s.GetString(0x1),
		Version:           s.GetString(0x2),
		SerialNumber:      s.GetString(0x3),
		AssetTag:          s.GetString(0x4),
		LocationInChassis: s.GetString(0x6),
		FeatureFlags:      FeatureFlags(s.GetByte(0x05)).String(),
		BoardType:         Type(s.GetByte(0x09)).String(),
	}
	return info.Map(), nil
}
