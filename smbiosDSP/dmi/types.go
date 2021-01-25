package dmi

type SmBiosStructureType uint8

const (
	BIOS SmBiosStructureType = iota /* 0 */
	System
	BaseBoard
	Chassis
	Processor
	MemoryController
	MemoryModule
	Cache
	PortConnector
	SystemSlots
	OnBoardDevices
	OEMStrings
	SystemConfigurationOptions
	BIOSLanguage
	GroupAssociations
	SystemEventLog
	PhysicalMemoryArray
	MemoryDevice
	Bit32MemoryError
	MemoryArrayMappedAddress
	MemoryDeviceMappedAddress
	BuiltInPointingDevice
	PortableBattery
	SystemReset
	HardwareSecurity
	SystemPowerControls
	VoltageProbe
	CoolingDevice
	TemperatureProbe
	ElectricalCurrentProbe
	OutOfBandRemoteAccess
	BootIntegrityServices
	SystemBoot
	Bit64MemoryError
	ManagementDevice
	ManagementDeviceComponent
	ManagementDeviceThresholdData
	MemoryChannel
	IPMIDevice
	PowerSupply
	AdditionalInformation
	OnBoardDevicesExtendedInformation
	ManagementControllerHostInterface
	TPMDevice                      /* 43 */
	Inactive   SmBiosStructureType = 126
	EndOfTable SmBiosStructureType = 127
)

func (b SmBiosStructureType) String() string {
	types := [...]string{
		"BIOS",
		"System",
		"BaseBoard",
		"Chassis",
		"Processor",
		"MemoryController",
		"MemoryModule",
		"Cache",
		"PortConnector",
		"SystemSlots",
		"OnBoardDevices",
		"OEMStrings",
		"SystemConfigurationOptions",
		"BIOSLanguage",
		"GroupAssociations",
		"SystemEventLog",
		"PhysicalMemoryArray",
		"MemoryDevice",
		"Bit32MemoryError",
		"MemoryArrayMappedAddress",
		"MemoryDeviceMappedAddress",
		"BuiltInPointingDevice",
		"PortableBattery",
		"SystemReset",
		"HardwareSecurity",
		"SystemPowerControls",
		"VoltageProbe",
		"CoolingDevice",
		"TemperatureProbe",
		"ElectricalCurrentProbe",
		"OutOfBandRemoteAccess",
		"BootIntegrityServices",
		"SystemBoot",
		"Bit64MemoryError",
		"ManagementDevice",
		"ManagementDeviceComponent",
		"ManagementDeviceThresholdData",
		"MemoryChannel",
		"IPMIDevice",
		"PowerSupply",
		"AdditionalInformation",
		"OnBoardDevicesExtendedInformation",
		"ManagementControllerHostInterface",
		"TPMDevice",
	}
	if b == 126 {
		return "Inactive"
	}
	if b == 127 {
		return "EndOfTable"
	}
	if b >= 128 {
		return "OEM-specific"
	}
	if b > 44 {
		return "Unknown"
	}
	return types[b]
}

var (
	SmBiosTypeInfo = make(map[string]SmBiosType)
)

type SmBiosType interface {
	Parse(s *SubSmBiosStructure) (interface{}, error)
}

func Register(name string, collect SmBiosType) {
	if collect == nil {
		panic("config: Register SmBiosType is nil")
	}
	if _, ok := SmBiosTypeInfo[name]; ok {
		panic("config: Register SmBiosType twice for adapter " + name)
	}
	SmBiosTypeInfo[name] = collect
}
