package riot

import (
	"gopher2600/hardware/memory/bus"
	"gopher2600/hardware/riot/timer"
	"strings"
)

// RIOT contains all the sub-components of the VCS RIOT sub-system
type RIOT struct {
	mem bus.ChipBus

	Timer *timer.Timer
}

// NewRIOT creates a RIOT, to be used in a VCS emulation
func NewRIOT(mem bus.ChipBus) *RIOT {
	riot := &RIOT{mem: mem}
	riot.Timer = timer.NewTimer(mem)

	return riot
}

func (riot RIOT) String() string {
	s := strings.Builder{}
	s.WriteString(riot.Timer.String())
	return s.String()
}

// ReadMemory checks for side effects to the RIOT sub-system
func (riot *RIOT) ReadMemory() {
	serviceMemory, data := riot.mem.ChipRead()
	if !serviceMemory {
		return
	}

	serviceMemory = riot.Timer.ServiceMemory(data)
	if !serviceMemory {
		return
	}

	// !!TODO: service other RIOT registers
}

// Step moves the state of the riot forward one video cycle
func (riot *RIOT) Step() {
	riot.ReadMemory()
	riot.Timer.Step()
}
