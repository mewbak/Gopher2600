package memory

import "fmt"

// VCSMemory presents a monolithic representation of system memory to the CPU -
// the CPU only ever access memory through an instance of this structure.
// Other parts of the system access ChipMemory directly
type VCSMemory struct {
	CPUBus
	memmap map[uint16]Area
	RIOT   *ChipMemory
	TIA    *ChipMemory
	PIA    *PIA
	Cart   *Cartridge
}

// NewVCSMemory is the preferred method of initialisation for VCSMemory
func NewVCSMemory() (*VCSMemory, error) {
	mem := new(VCSMemory)
	if mem == nil {
		return nil, fmt.Errorf("can't allocate memory for VCS")
	}

	mem.memmap = make(map[uint16]Area)
	mem.RIOT = NewRIOT()
	mem.TIA = NewTIA()
	mem.PIA = NewPIA()
	mem.Cart = NewCart()

	if mem.memmap == nil || mem.RIOT == nil || mem.TIA == nil || mem.PIA == nil || mem.Cart == nil {
		return nil, fmt.Errorf("can't allocate memory for VCS")
	}

	// create the memory map; each address in the memory map points to the
	// memory area it resides in. we only record 'primary' addresses; all
	// addresses should be  passed through the MapAddress() function in order
	// to iron out any mirrors

	var i uint16

	for i = mem.TIA.origin; i <= mem.TIA.memtop; i++ {
		mem.memmap[i] = mem.TIA
	}

	for i = mem.PIA.origin; i <= mem.PIA.memtop; i++ {
		mem.memmap[i] = mem.PIA
	}

	for i = mem.RIOT.origin; i <= mem.RIOT.memtop; i++ {
		mem.memmap[i] = mem.RIOT
	}

	for i = mem.Cart.origin; i <= mem.Cart.memtop; i++ {
		mem.memmap[i] = mem.Cart
	}

	return mem, nil
}

func (mem VCSMemory) String() string {
	return mem.MemoryMap()
}

// MapAddress translates the quoted address from mirror space to primary space.
// Generally, all access to the different memory areas should be passed through
// this function. Any other information about an address can be accessed
// through mem.memmap[mappedAddress]
func (mem VCSMemory) MapAddress(address uint16) uint16 {
	// note that the order of these filters is important

	// cartridge addresses
	if address&mem.Cart.origin == mem.Cart.origin {
		return address & mem.Cart.memtop
	}

	// RIOT addresses
	if address&mem.RIOT.origin == mem.RIOT.origin {
		return address & mem.RIOT.memtop
	}

	// PIA RAM addresses
	if address&mem.PIA.origin == mem.PIA.origin {
		return address & mem.PIA.memtop
	}

	// everything else is in TIA space
	return address & mem.TIA.memtop
}

// Clear is an implementation of CPUBus.Clear
func (mem *VCSMemory) Clear() {
	mem.RIOT.Clear()
	mem.TIA.Clear()
	mem.PIA.Clear()
	mem.Cart.Clear()
}

// Implementation of CPUBus.Read
func (mem VCSMemory) Read(address uint16) (uint8, error) {
	ma := mem.MapAddress(address)
	area, present := mem.memmap[ma]
	if !present {
		return 0, fmt.Errorf("%04x not mapped correctly", address)
	}
	return area.(CPUBus).Read(ma)
}

// Implementation of CPUBus.Write
func (mem *VCSMemory) Write(address uint16, data uint8) error {
	ma := mem.MapAddress(address)
	area, present := mem.memmap[ma]
	if !present {
		return fmt.Errorf("%04x not mapped correctly", address)
	}
	return area.(CPUBus).Write(ma, data)
}
