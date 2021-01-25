package smbiosDSP

import (
	"github.com/shhnwangjian/ops-warren/smbiosDSP/dmi"
	"github.com/shhnwangjian/ops-warren/smbiosDSP/typeStructure/bios"
)

type Decoder struct {
	eps                    *dmi.EntryPoint
	bios                   []*dmi.SubSmBiosStructure
	system                 []*dmi.SubSmBiosStructure
	baseBoard              []*dmi.SubSmBiosStructure
	chassis                []*dmi.SubSmBiosStructure
	onBoardDevices         []*dmi.SubSmBiosStructure
	onBoardExtendedDevices []*dmi.SubSmBiosStructure
	portConnector          []*dmi.SubSmBiosStructure
	processor              []*dmi.SubSmBiosStructure
	cache                  []*dmi.SubSmBiosStructure
	physicalMemoryArray    []*dmi.SubSmBiosStructure
	memoryDevice           []*dmi.SubSmBiosStructure
	systemSlots            []*dmi.SubSmBiosStructure
	portableBattery        []*dmi.SubSmBiosStructure
}

func New() (*Decoder, error) {
	eps, ss, err := dmi.ReadSmBios()
	if err != nil {
		return nil, err
	}

	d := new(Decoder)
	d.eps = eps

	for i := range ss {
		switch dmi.SmBiosStructureType(ss[i].Header.Type) {
		case dmi.BIOS:
			d.bios = append(d.bios, ss[i])
		case dmi.System:
			d.system = append(d.system, ss[i])
		case dmi.BaseBoard:
			d.baseBoard = append(d.baseBoard, ss[i])
		case dmi.Chassis:
			d.chassis = append(d.chassis, ss[i])
		case dmi.OnBoardDevices:
			d.onBoardDevices = append(d.onBoardDevices, ss[i])
		case dmi.OnBoardDevicesExtendedInformation:
			d.onBoardExtendedDevices = append(d.onBoardExtendedDevices, ss[i])
		case dmi.PortConnector:
			d.portConnector = append(d.portConnector, ss[i])
		case dmi.Processor:
			d.processor = append(d.processor, ss[i])
		case dmi.Cache:
			d.cache = append(d.cache, ss[i])
		case dmi.PhysicalMemoryArray:
			d.physicalMemoryArray = append(d.physicalMemoryArray, ss[i])
		case dmi.MemoryDevice:
			d.memoryDevice = append(d.memoryDevice, ss[i])
		case dmi.SystemSlots:
			d.systemSlots = append(d.systemSlots, ss[i])
		case dmi.PortableBattery:
			d.portableBattery = append(d.portableBattery, ss[i])
		default:
		}
	}
	return d, nil
}

// BIOS 解析bios信息
func (d *Decoder) BIOS() (infos []map[string]string, err error) {
	infos = make([]map[string]string, 0)
	for i := range d.bios {
		info, err := bios.Parse(d.bios[i])
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}
