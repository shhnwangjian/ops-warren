package battery

import (
	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

type DeviceChemistry byte

const (
	Other DeviceChemistry = 1 << iota
	Unknown
	LeadAcid
	NickelCadmium
	NickelMetalHydride
	LithiumIon
	ZincAir
	LithiumPolymer
)

func (d DeviceChemistry) String() string {
	factors := [...]string{
		"Other",
		"Unknown",
		"Lead Acid",
		"Nickel Cadmium",
		"Nickel metal hydride",
		"Lithium-ion",
		"Zinc air",
		"Lithium Polymer",
	}
	return factors[d]
}

type Information struct {
	dmi.Header
	Location                 string `json:"location,omitempty"`
	Manufacturer             string `json:"manufacturer,omitempty"`
	ManufacturerDate         string `json:"manufacturer_date,omitempty"`
	SerialNumber             string `json:"serial_number,omitempty"`
	DeviceName               string `json:"device_name,omitempty"`
	DeviceChemistry          string `json:"device_chemistry,omitempty"`
	DesignCapacity           uint16 `json:"design_chemistry,omitempty"`
	DesignVoltage            uint16 `json:"design_voltage,omitempty"`
	SBDSVersionNumber        string `json:"sbds_version_number,omitempty"`
	MaximumError             uint16 `json:"maximum_error,omitempty"`
	SBDSSerialNumber         uint16 `json:"sbds_serial_number,omitempty"`
	SBDSManufactureDate      uint16 `json:"sbds_manufacture_date,omitempty"`
	SBDSDeviceChemistry      string `json:"sbds_device_chemistry,omitempty"`
	DesignCapacityMultiplier uint16 `json:"design_capacity_multiplier,omitempty"`
	OEMSpecific              uint16 `json:"oem_specific,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Location":            b.Location,
		"Manufacturer":        b.Manufacturer,
		"ManufacturerDate":    b.ManufacturerDate,
		"SerialNumber":        b.SerialNumber,
		"DeviceName":          b.DeviceName,
		"DeviceChemistry":     b.DeviceChemistry,
		"SBDSVersionNumber":   b.SBDSVersionNumber,
		"SBDSDeviceChemistry": b.SBDSDeviceChemistry,
	}
}

// Parse 解析电池信息
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		Header:              s.Header,
		Location:            s.GetString(0x00),
		Manufacturer:        s.GetString(0x01),
		ManufacturerDate:    s.GetString(0x02),
		SerialNumber:        s.GetString(0x03),
		DeviceName:          s.GetString(0x04),
		DeviceChemistry:     DeviceChemistry(s.GetByte(0x05)).String(),
		SBDSVersionNumber:   s.GetString(0x0a),
		SBDSDeviceChemistry: s.GetString(0x0e),
	}
	designCapacity, err := s.U16(0x06, 0x08)
	if err == nil {
		info.DesignCapacity = designCapacity
	}
	designVoltage, err := s.U16(0x08, 0x0a)
	if err == nil {
		info.DesignVoltage = designVoltage
	}
	maximumError, err := s.U16(0x0b, 0x0c)
	if err == nil {
		info.MaximumError = maximumError
	}
	sBDSManufactureDate, err := s.U16(0x0c, 0x0e)
	if err == nil {
		info.SBDSManufactureDate = sBDSManufactureDate
	}
	designCapacityMultiplier, err := s.U16(0x0f, 0x10)
	if err == nil {
		info.DesignCapacityMultiplier = designCapacityMultiplier
	}
	oEMSpecific, err := s.U16(0x10, 0x11)
	if err == nil {
		info.OEMSpecific = oEMSpecific
	}
	return info.Map(), nil

}
