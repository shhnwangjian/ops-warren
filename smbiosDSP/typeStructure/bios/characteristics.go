package bios

import "fmt"

// Characteristics bios字符集
type Characteristics uint64

// BIOS Characteristics
const (
	Reserved0 Characteristics = 1 << iota
	Reserved1
	Unknown
	NotSupported
	ISASupported
	MCASupported
	EISASupported
	PCISupported
	PCMCIASupported
	PlugPlaySupported
	APMSupported
	Upgradeable
	ShadowingIsAllowed
	VLVESASupported
	ESCDSupported
	BootFromCDSupported
	SelectableBootSupported
	BIOSROMIsSocketed
	BootFromPCMCIASupported
	EDDSupported
	JPFloppyNECSupported
	JPFloppyToshibaSupported
	Floppy525_360KBSupported
	Floppy525_1_2MBSupported
	Floppy35_720KBSupported
	Floppy35_2_88MBSupported
	PrintScreenSupported
	Keyboard8042Supported
	SerialSupported
	PrinterSupported
	CGAMonoSupported
	NECPC98
	//Bit32:47 Reserved for BIOS vendor
	//Bit47:63 Reserved for system vendor
)

func (b Characteristics) String() string {
	var s string
	chars := [...]string{
		"BIOS characteristics not supported", /* 3 */
		"ISA is supported",
		"MCA is supported",
		"EISA is supported",
		"PCI is supported",
		"PC Card (PCMCIA) is supported",
		"PNP is supported",
		"APM is supported",
		"BIOS is upgradeable",
		"BIOS shadowing is allowed",
		"VLB is supported",
		"ESCD support is available",
		"Boot from CD is supported",
		"Selectable boot is supported",
		"BIOS ROM is socketed",
		"Boot from PC Card (PCMCIA) is supported",
		"EDD is supported",
		"Japanese floppy for NEC 9800 1.2 MB is supported (int 13h)",
		"Japanese floppy for Toshiba 1.2 MB is supported (int 13h)",
		"5.25\"/360 kB floppy services are supported (int 13h)",
		"5.25\"/1.2 MB floppy services are supported (int 13h)",
		"3.5\"/720 kB floppy services are supported (int 13h)",
		"3.5\"/2.88 MB floppy services are supported (int 13h)",
		"Print screen service is supported (int 5h)",
		"8042 keyboard services are supported (int 9h)",
		"Serial services are supported (int 14h)",
		"Printer services are supported (int 17h)",
		"CGA/mono video services are supported (int 10h)",
		"NEC PC-98", /* 31 */
	}
	for i := uint32(4); i < 32; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i-3]
		}
	}
	return s
}

// Ext1 ext1
type Ext1 byte

// BIOS Characteristics Extension Bytes(Ext1)
// Byte 1
const (
	Ext1ACPISupported Ext1 = 1 << iota
	Ext1USBLegacySupported
	Ext1AGPSupported
	Ext1I2OBootSupported
	Ext1LS120SupperDiskBootSupported
	Ext1ATAPIZIPDriveBootSupported
	Ext11394BootSupported
	Ext1SmartBatterySupported
)

func (b Ext1) String() string {
	var s string
	chars := [...]string{
		"ACPI is supported", /* 0 */
		"USB legacy is supported",
		"AGP is supported",
		"I2O boot is supported",
		"LS-120 boot is supported",
		"ATAPI Zip drive boot is supported",
		"IEEE 1394 boot is supported",
		"Smart battery is supported", /* 7 */
	}

	for i := uint32(0); i < 7; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i]
		}
	}
	return s
}

// Ext2 ext2
type Ext2 byte

// BIOS Characteristics Extension Bytes(Ext2)
// Byte 2
const (
	Ext2BIOSBootSpecSupported Ext2 = 1 << iota
	Ext2FuncKeyInitiatedNetworkBootSupported
	Ext2EnableTargetedContentDistribution
	Ext2UEFISpecSupported
	Ext2VirtualMachine
	// Bits 5:7 Reserved for future assignment
)

func (b Ext2) String() string {
	var s string
	chars := [...]string{
		"BIOS boot specification is supported", /* 0 */
		"Function key-initiated network boot is supported",
		"Targeted content distribution is supported",
		"UEFI is supported",
		"System is a virtual machine", /* 4 */
	}

	for i := uint32(0); i < 5; i++ {
		if b&(1<<i) != 0 {
			s += "\n\t\t" + chars[i]
		}
	}
	return s
}

// RuntimeSize
type RuntimeSize uint

func (b RuntimeSize) String() string {
	if (b & 0x3FF) > 0 {
		return fmt.Sprintf("%d Bytes", b)
	}
	return fmt.Sprintf("%d kB", b>>10)
}

// RomSize
type RomSize byte

func (b RomSize) String() string {
	return fmt.Sprintf("%d kB", b)
}
