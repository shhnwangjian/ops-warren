package chassis

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

const (
	OUT_OF_SPEC = "<OUT OF SPEC>"
)

// Type Chssis 类型
type Type byte

const (
	Other Type = 1 + iota
	Unknown
	Desktop
	LowProfileDesktop
	PizzaBox
	MiniTower
	Tower
	Portable
	Laptop
	Notebook
	HandHeld
	DockingStation
	AllinOne
	SubNotebook
	SpaceSaving
	LunchBox
	MainServerChassis
	ExpansionChassis
	SubChassis
	BusExpansionChassis
	PeripheralChassis
	RAIDChassis
	RackMountChassis
	SealedcasePC
	MultiSystem
	CompactPCI
	AdvancedTCA
	Blade
	BladeEnclosure
)

func (c Type) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"Desktop",
		"LowProfileDesktop",
		"PizzaBox",
		"MiniTower",
		"Tower",
		"Portable",
		"Laptop",
		"Notebook",
		"HandHeld",
		"Docking Station",
		"AllinOne",
		"SubNotebook",
		"SpaceSaving",
		"LunchBox",
		"MainServer Chassis",
		"Expansion Chassis",
		"Sub Chassis",
		"BusExpansion Chassis",
		"Peripheral Chassis",
		"RAID Chassis",
		"Rack Mount Chassis",
		"SealedcasePC",
		"MultiSystem",
		"CompactPCI",
		"AdvancedTCA",
		"Blade",
		"BladeEnclosure",
	}
	c &= 0x7F
	if c >= 0x01 && c < 0x1D {
		return types[c-1]
	}
	return OUT_OF_SPEC
}

// Lock todo
type Lock byte

func (c Lock) String() string {
	locks := [...]string{
		"Not Present", /* 0x00 */
		"Present",     /* 0x01 */
	}
	return locks[c]
}

// State todo
type State byte

const (
	StateOther State = 1 + iota
	StateUnknown
	StateSafe
	StateWarning
	StateCritical
	StateNonRecoverable
)

func (c State) String() string {
	states := [...]string{
		"Other",
		"Unknown",
		"Safe",
		"Warning",
		"Critical",
		"NonRecoverable",
	}
	return states[c-1]
}

type ContainedElementType byte

type ContainedElements struct {
	Type    ContainedElementType
	Minimum byte
	Maximum byte
}

type SecurityStatus byte

const (
	ChassisSecurityStatusOther SecurityStatus = 1 + iota
	ChassisSecurityStatusUnknown
	ChassisSecurityStatusNone
	ChassisSecurityStatusExternalInterfaceLockedOut
	ChassisSecurityStatusExternalInterfaceEnabled
)

func (s SecurityStatus) String() string {
	status := [...]string{
		"Other",
		"Unknown",
		"None",
		"ExternalInterfaceLockedOut",
		"ExternalInterfaceEnabled",
	}
	return status[s-1]
}

// Height 高度
type Height byte

func (s Height) String() string {
	if s == 0 {
		return fmt.Sprintf("Not Specified")
	}

	return fmt.Sprintf("%d U", s)
}

// Information 底座信息
type Information struct {
	dmi.Header
	Manufacturer                 string `json:"manufacturer,omitempty"`
	Type                         string `json:"type,omitempty"`
	Lock                         string `json:"lock,omitempty"`
	Version                      string `json:"version,omitempty"`
	AssetTag                     string `json:"asset_tag,omitempty"`
	SerialNumber                 string `json:"serial_number,omitempty"`
	BootUpState                  string `json:"boot_up_state,omitempty"`
	PowerSupplyState             string `json:"power_supply_state,omitempty"`
	ThermalState                 string `json:"thermal_state,omitempty"`
	SecurityStatus               string `json:"security_status,omitempty"`
	OEMdefined                   uint16 `json:"oe_mdefined,omitempty"`
	Height                       string `json:"height,omitempty"`
	NumberOfPowerCords           byte   `json:"number_of_power_cords,omitempty"`
	ContainedElementCount        byte   `json:"contained_element_count,omitempty"`
	ContainedElementRecordLength byte   `json:"contained_element_record_length,omitempty"`
	//ContainedElements            ContainedElements `json:"contained_elements,omitempty"`
	SKUNumber string `json:"sku_number,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Manufacturer":     b.Manufacturer,
		"Type":             b.Type,
		"Version":          b.Version,
		"SerialNumber":     b.SerialNumber,
		"AssetTag":         b.AssetTag,
		"Lock":             b.Lock,
		"BootUpState":      b.BootUpState,
		"PowerSupplyState": b.PowerSupplyState,
		"ThermalState":     b.ThermalState,
		"SecurityStatus":   b.SecurityStatus,
		"Height":           b.Height,
		"SKUNumber":        b.SKUNumber,
	}
}

// Parse 解析底座信息
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		Manufacturer:                 s.GetString(0x0),
		Type:                         Type(s.GetByte(0x01) & 127).String(),
		Lock:                         Lock(s.GetByte(0x01) >> 7).String(),
		Version:                      s.GetString(0x2),
		SerialNumber:                 s.GetString(0x3),
		AssetTag:                     s.GetString(0x4),
		BootUpState:                  State(s.GetByte(0x05)).String(),
		PowerSupplyState:             State(s.GetByte(0x06)).String(),
		ThermalState:                 State(s.GetByte(0x07)).String(),
		SecurityStatus:               SecurityStatus(s.GetByte(0x08)).String(),
		Height:                       Height(s.GetByte(0xd)).String(),
		NumberOfPowerCords:           s.GetByte(0xe),
		ContainedElementCount:        s.GetByte(0xf),
		ContainedElementRecordLength: s.GetByte(0x10),
		// TODO: 7.4.4
		//ci.ContainedElements:
		SKUNumber: s.GetString(4),
	}
	oemDefined, err := s.U16(0x09, 0xd)
	if err == nil {
		info.OEMdefined = oemDefined
	}
	return info.Map(), nil
}
