package debugger

// Help contains the help text for the debugger's top level commands
var Help = map[string]string{
	cmdBall:          "Display the current state of the ball sprite",
	cmdBreak:         "Cause emulator to halt when conditions are met",
	cmdCPU:           "Display the current state of the CPU",
	cmdCartridge:     "Display information about the current cartridge",
	cmdClear:         "Clear all entries in BREAKS and TRAPS",
	cmdDebuggerState: "Display summary of debugger options",
	cmdDigest:        "Return the cryptographic hash of the current screen",
	cmdDisassembly:   "Print the full cartridge disassembly",
	cmdDisplay:       "Display the TV image",
	cmdDrop:          "Drop a specific BREAK or TRAP conditin, using the number of the condition reported by LIST",
	cmdGrep:          "Simple string search (case insensitive) of the disassembly",
	cmdHelp:          "Lists commands and provides help for individual debugger commands",
	cmdHexLoad:       "Modify a sequence of memory addresses. Starting address must be numeric.",
	cmdInsert:        "Insert cartridge into emulation (from file)",
	cmdLast:          "Prints the result of the last cpu/video cycle",
	cmdList:          "List current entries for BREAKS and TRAPS",
	cmdMemMap:        "Display high-level VCS memory map",
	cmdMissile:       "Display the current state of the missile 0/1 sprite",
	cmdOnHalt:        "Commands to run whenever emulation is halted (separate commands with comma)",
	cmdOnStep:        "Commands to run whenever emulation steps forward an cpu/video cycle (separate commands with comma)",
	cmdPeek:          "Inspect an individual memory address",
	cmdPlayer:        "Display the current state of the player 0/1 sprite",
	cmdPlayfield:     "Display the current playfield data",
	cmdPoke:          "Modify an individual memory address",
	cmdQuit:          "Exits the emulator",
	cmdRAM:           "Display the current contents of PIA RAM",
	cmdRIOT:          "Display the current state of the RIOT",
	cmdReset:         "Reset the emulation to its initial state",
	cmdRun:           "Run emulator until next halt state",
	cmdScript:        "Run commands from specified file or record commands to a file",
	cmdStep:          "Step forward one step. Optional argument sets the amount to step by (eg. frame, scanline, etc.)",
	cmdGranularity:   "Change method of stepping: CPU or VIDEO",
	cmdStick:         "Emulate a joystick input for Player 0 or Player 1",
	cmdSymbol:        "Search for the address label symbol in disassembly. returns address",
	cmdTIA:           "Display current state of the TIA",
	cmdTV:            "Display the current TV state",
	cmdTerse:         "Use terse format when displaying machine information",
	cmdTrap:          "Cause emulator to halt when specified machine component is touched",
	cmdVerbose:       "Use verbose format when displaying machine information",
	cmdVerbosity:     "Display which format is used when displaying machine information (see TERSE and VERBOSE commands)",
	cmdWatch:         "Watch a memory address for activity",
}
