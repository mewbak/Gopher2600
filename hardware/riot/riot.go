package riot

import (
	"fmt"
	"gopher2600/hardware/memory"
	"gopher2600/hardware/memory/vcssymbols"
)

// RIOT contains all the sub-components of the VCS RIOT sub-system
type RIOT struct {
	mem memory.ChipBus

	// timerRegister is the name of the currently selected RIOT timer. used as a
	// label in MachineInfo()
	timerRegister string

	// timerInterval indicates how often (in CPU cycles) the timer value
	// descreases
	timerInterval int

	// timerINTIMvalue is the current timer value and is reflected in the INTIM
	// register (RIOT memory)
	timerINTIMvalue uint8

	// timerCycles is the number of CPU cycles remainng before INTIM is decreased
	// when a new time is started, timerCycles is always set to two (decrease
	// occurs almost immediately) and thereafter set to the selected
	// timerInterval
	//
	// the initial reset value is 2 because the first decrease of INTIM occurs on
	// the *next* machine cycle - timerCycles will be reduced to 1 on the same
	// machine cycle it is set to 2, and to 0 on the *next* cycle. phew.
	timerCycles int
}

// NewRIOT creates a RIOT, to be used in a VCS emulation
func NewRIOT(mem memory.ChipBus) *RIOT {
	riot := new(RIOT)
	riot.timerRegister = "no timer"
	riot.mem = mem
	return riot
}

// MachineInfoTerse returns the RIOT information in terse format
func (riot RIOT) MachineInfoTerse() string {
	return fmt.Sprintf("INTIM=%d clks=%d (%s)", riot.timerINTIMvalue, riot.timerCycles, riot.timerRegister)
}

// MachineInfo returns the RIOT information in verbose format
func (riot RIOT) MachineInfo() string {
	return fmt.Sprintf("%s\nINTIM: %d (%02x)\nINTIM clocks = %d (%02x)", riot.timerRegister, riot.timerINTIMvalue, riot.timerINTIMvalue, riot.timerCycles, riot.timerCycles)
}

// map String to MachineInfo
func (riot RIOT) String() string {
	return riot.MachineInfo()
}

// ReadRIOTMemory checks for side effects to the RIOT sub-system
func (riot *RIOT) ReadRIOTMemory() {
	service, register, value := riot.mem.ChipRead()
	if service {
		switch register {
		case "TIM1T":
			riot.timerRegister = register
			riot.timerInterval = 1
			riot.timerINTIMvalue = value
			riot.timerCycles = 2
		case "TIM8T":
			riot.timerRegister = register
			riot.timerInterval = 8
			riot.timerINTIMvalue = value
			riot.timerCycles = 2
		case "TIM64T":
			riot.timerRegister = register
			riot.timerInterval = 64
			riot.timerINTIMvalue = value
			riot.timerCycles = 2
		case "TIM1024":
			riot.timerRegister = register
			riot.timerInterval = 1024
			riot.timerINTIMvalue = value
			riot.timerCycles = 2

			// TODO: handle other RIOT registers
		}

		// write value to INTIM straight-away
		riot.mem.ChipWrite(vcssymbols.INTIM, uint8(riot.timerINTIMvalue))
	}
}

// Step moves the state of the riot forward one video cycle
func (riot *RIOT) Step() {
	// some documentation (Atari 2600 Specifications.htm) claims that if INTIM is
	// *read* then the decrement reverts to once per timer interval. this won't
	// have any discernable effect unless the timer interval has been flipped to
	// 1 when INTIM cycles back to 255
	if riot.mem.LastReadRegister() == "INTIM" {
		switch riot.timerRegister {
		case "TIM1T":
			riot.timerInterval = 1
		case "TIM8T":
			riot.timerInterval = 8
		case "TIM64T":
			riot.timerInterval = 64
		case "TIM1024":
			riot.timerInterval = 1024
		}
	}

	if riot.timerRegister != "no timer" {
		riot.timerCycles--
		if riot.timerCycles == 0 {
			if riot.timerINTIMvalue == 0 {
				// reset INTIM value
				riot.timerINTIMvalue = 255

				// because INTIM value has cycled we flip timer interval to 1
				riot.timerInterval = 1
			} else {
				riot.timerINTIMvalue--
			}
			riot.mem.ChipWrite(vcssymbols.INTIM, riot.timerINTIMvalue)
			riot.timerCycles = riot.timerInterval
		}
	}
}
