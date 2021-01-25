package dmi

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/digitalocean/go-smbios/smbios"
)

// https://www.dmtf.org/sites/default/files/standards/documents/DSP0134_3.1.1.pdf

// A Header is a Structure's header.
type Header struct {
	Type   uint8
	Length uint8
	Handle uint16
}

type SubSmBiosStructure struct {
	Header    Header
	Formatted []byte
	Strings   []string
}

// NewEntryPoint
func NewEntryPoint(major, minor, rev, addr, size int) *EntryPoint {
	return &EntryPoint{
		Major:    major,
		Minor:    minor,
		Revision: rev,
		Address:  addr,
		Size:     size,
	}
}

// EntryPoint EPS
type EntryPoint struct {
	Address  int `json:"address,omitempty"`
	Size     int `json:"size,omitempty"`
	Major    int `json:"major,omitempty"`
	Minor    int `json:"minor,omitempty"`
	Revision int `json:"revision,omitempty"`
}

// ReadSmBios 读取smbios结构数据
func ReadSmBios() (*EntryPoint, []*SubSmBiosStructure, error) {
	// Find SMBIOS data in operating system-specific location.
	rc, ep, err := smbios.Stream()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stream: %v", err)
	}
	// Be sure to close the stream!
	defer rc.Close()

	// Decode SMBIOS structures from the stream.
	d := smbios.NewDecoder(rc)
	ss, err := d.Decode()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode structures: %v", err)
	}

	// Determine SMBIOS version and table location from entry point.
	major, minor, rev := ep.Version()
	addr, size := ep.Table()

	eps := NewEntryPoint(major, minor, rev, addr, size)
	data := allSubSmBiosStructure(ss)

	return eps, data, nil
}

func allSubSmBiosStructure(ss []*smbios.Structure) []*SubSmBiosStructure {
	data := make([]*SubSmBiosStructure, 0, len(ss))
	for _, s := range ss {
		ns := SubSmBiosStructure{
			Header: Header{
				Type:   s.Header.Type,
				Length: s.Header.Length,
				Handle: s.Header.Handle,
			},
			Formatted: s.Formatted,
			Strings:   s.Strings,
		}
		data = append(data, &ns)
	}
	return data
}

func (s *SubSmBiosStructure) GetString(offset int) string {
	if offset > len(s.Formatted)-1 {
		return "Unknown"
	}

	index := s.Formatted[offset]

	if index == 0 {
		return "Unknown"
	}
	return s.Strings[index-1]
}

func (s *SubSmBiosStructure) GetBytes(start, end int) []byte {
	if s.IsOverStep(end) {
		return []byte{}
	}

	return s.Formatted[start:end]
}

func (s *SubSmBiosStructure) GetByte(index int) byte {
	if s.IsOverStep(index) {
		return 0
	}
	return s.Formatted[index]
}

func (s *SubSmBiosStructure) DataLength() int {
	return len(s.Formatted)
}

func (s *SubSmBiosStructure) Type() uint8 {
	return s.Header.Type
}

func (s *SubSmBiosStructure) IsOverStep(index int) bool {
	return index+1 > s.DataLength()
}

func (s *SubSmBiosStructure) U16(start, end int) (uint16, error) {
	if s.IsOverStep(end) {
		return 0, nil
	}

	return U16(s.Formatted[start:end])
}

func (s *SubSmBiosStructure) U32(start, end int) (uint32, error) {
	if s.IsOverStep(end) {
		return 0, nil
	}

	return U32(s.Formatted[start:end])
}

func (s *SubSmBiosStructure) U64(start, end int) (uint64, error) {
	if s.IsOverStep(end) {
		return 0, nil
	}

	return U64(s.Formatted[start:end])
}

func U16(data []byte) (uint16, error) {
	var u uint16
	err := binary.Read(bytes.NewBuffer(data[0:2]), binary.LittleEndian, &u)
	return u, err
}

func U32(data []byte) (uint32, error) {
	var u uint32
	err := binary.Read(bytes.NewBuffer(data[0:4]), binary.LittleEndian, &u)
	return u, err
}

func U64(data []byte) (uint64, error) {
	var u uint64
	err := binary.Read(bytes.NewBuffer(data[0:8]), binary.LittleEndian, &u)
	return u, err
}
