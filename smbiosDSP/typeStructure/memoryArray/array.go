package memoryArray

import (
	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
)

type PhysicalMemoryArrayLocation byte

const (
	PhysicalMemoryArrayLocationOther PhysicalMemoryArrayLocation = 1 + iota
	PhysicalMemoryArrayLocationUnknown
	PhysicalMemoryArrayLocationSystemboardormotherboard
	PhysicalMemoryArrayLocationISAadd_oncard
	PhysicalMemoryArrayLocationEISAadd_oncard
	PhysicalMemoryArrayLocationPCIadd_oncard
	PhysicalMemoryArrayLocationMCAadd_oncard
	PhysicalMemoryArrayLocationPCMCIAadd_oncard
	PhysicalMemoryArrayLocationProprietaryadd_oncard
	PhysicalMemoryArrayLocationNuBus
	PhysicalMemoryArrayLocationPC_98C20add_oncard
	PhysicalMemoryArrayLocationPC_98C24add_oncard
	PhysicalMemoryArrayLocationPC_98Eadd_oncard
	PhysicalMemoryArrayLocationPC_98Localbusadd_oncard
)

func (p PhysicalMemoryArrayLocation) String() string {
	locations := [...]string{
		"Other",
		"Unknown",
		"System board or motherboard",
		"ISA add-on card",
		"EISA add-on card",
		"PCI add-on card",
		"MCA add-on card",
		"PCMCIA add-on card",
		"Proprietary add-on card",
		"NuBus",
		"PC-98/C20 add-on card",
		"PC-98/C24 add-on card",
		"PC-98/E add-on card",
		"PC-98/Local bus add-on card",
	}
	return locations[p-1]
}

type PhysicalMemoryArrayUse byte

const (
	PhysicalMemoryArrayUseOther PhysicalMemoryArrayUse = 1 + iota
	PhysicalMemoryArrayUseUnknown
	PhysicalMemoryArrayUseSystemmemory
	PhysicalMemoryArrayUseVideomemory
	PhysicalMemoryArrayUseFlashmemory
	PhysicalMemoryArrayUseNon_volatileRAM
	PhysicalMemoryArrayUseCachememory
)

func (p PhysicalMemoryArrayUse) String() string {
	uses := [...]string{
		"Other",
		"Unknown",
		"System memory",
		"Video memory",
		"Flash memory",
		"Non-volatile RAM",
		"Cache memory",
	}
	return uses[p-1]
}

type PhysicalMemoryArrayErrorCorrection byte

const (
	PhysicalMemoryArrayErrorCorrectionOther PhysicalMemoryArrayErrorCorrection = 1 + iota
	PhysicalMemoryArrayErrorCorrectionUnknown
	PhysicalMemoryArrayErrorCorrectionNone
	PhysicalMemoryArrayErrorCorrectionParity
	PhysicalMemoryArrayErrorCorrectionSingle_bitECC
	PhysicalMemoryArrayErrorCorrectionMulti_bitECC
	PhysicalMemoryArrayErrorCorrectionCRC
)

func (p PhysicalMemoryArrayErrorCorrection) String() string {
	types := [...]string{
		"Other",
		"Unknown",
		"None",
		"Parity",
		"Single-bit ECC",
		"Multi-bit ECC",
		"CRC",
	}
	return types[p-1]
}

// PhysicalMemoryArray todo
type Information struct {
	dmi.Header
	Location                string `json:"location,omitempty"`
	Use                     string `json:"use,omitempty"`
	ErrorCorrection         string `json:"error_correction,omitempty"`
	MaximumCapacity         uint32 `json:"maximum_capacity,omitempty"`
	ErrorInformationHandle  uint16 `json:"error_information_handle,omitempty"`
	NumberOfMemoryDevices   uint16 `json:"number_of_memory_devices,omitempty"`
	ExtendedMaximumCapacity uint64 `json:"extended_maximum_capacity,omitempty"`
}

func (b Information) Map() map[string]string {
	return map[string]string{
		"Location":        b.Location,
		"Use":             b.Use,
		"ErrorCorrection": b.ErrorCorrection,
	}
}

// Parse Physical Memory Array
func Parse(s *dmi.SubSmBiosStructure) (map[string]string, error) {
	info := &Information{
		Header:          s.Header,
		Location:        PhysicalMemoryArrayLocation(s.GetByte(0x00)).String(),
		Use:             PhysicalMemoryArrayUse(s.GetByte(0x01)).String(),
		ErrorCorrection: PhysicalMemoryArrayErrorCorrection(s.GetByte(0x02)).String(),
	}
	maximumCapacity, err := s.U32(0x03, 0x07)
	if err == nil {
		info.MaximumCapacity = maximumCapacity
	}
	errorInformationHandle, err := s.U16(0x07, 0x09)
	if err == nil {
		info.ErrorInformationHandle = errorInformationHandle
	}
	numberOfMemoryDevices, err := s.U16(0x09, 0xb)
	if err == nil {
		info.NumberOfMemoryDevices = numberOfMemoryDevices
	}
	extendedMaximumCapacity, err := s.U64(0xb, 0x13)
	if err == nil {
		info.ExtendedMaximumCapacity = extendedMaximumCapacity
	}
	return info.Map(), nil
}
