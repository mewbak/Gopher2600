package debugger

import (
	"gopher2600/cartridgeloader"
	"gopher2600/debugger/commandline"
	"gopher2600/debugger/console"
	"gopher2600/debugger/reflection"
	"gopher2600/debugger/script"
	"gopher2600/disassembly"
	"gopher2600/errors"
	"gopher2600/gui"
	"gopher2600/hardware"
	"gopher2600/screendigest"
	"gopher2600/setup"
	"gopher2600/symbols"
	"gopher2600/television"
	"os"
	"os/signal"
	"strings"
)

const defaultOnHalt = "CPU; TV"
const defaultOnStep = "LAST"
const onEmptyInput = "STEP"

// Debugger is the basic debugging frontend for the emulation
type Debugger struct {
	vcs    *hardware.VCS
	disasm *disassembly.Disassembly

	// gui/tv
	tv     television.Television
	scr    gui.GUI
	digest *screendigest.SHA1

	// interface to the vcs memory with additional debugging functions
	// -- access to vcs memory from the debugger (eg. peeking and poking) is
	// most fruitfully performed through this structure
	dbgmem *memoryDebug

	// reflection is used to provideo additional information about the
	// emulation. it is inherently slow so should be deactivated if not
	// required
	relfectMonitor *reflection.Monitor

	// halt conditions
	breakpoints *breakpoints
	traps       *traps
	watches     *watches

	// single-fire step traps. these are used for the STEP command, allowing
	// things like "STEP FRAME".
	stepTraps *traps

	// commandOnHalt says whether an sequence of commands should run automatically
	// when emulation halts. commandOnHaltPrev is the stored command sequence
	// used when ONHALT is called with no arguments
	// halt is a breakpoint or user intervention (ie. ctrl-c)
	commandOnHalt       string
	commandOnHaltStored string

	// similarly, commandOnStep is the sequence of commands to run afer every
	// cpu/video cycle
	commandOnStep       string
	commandOnStepStored string

	// whether to display the triggering of a known CPU bug. these are bugs
	// that are known about in the emulated hardware but which might catch an
	// unwary programmer by surprise
	reportCPUBugs bool

	// granularity of single stepping - every cpu instruction or every video cycle
	// -- also affects when emulation will halt on breaks, traps and watches.
	// if inputeveryvideocycle is true then the halt may occur mid-cpu-cycle
	inputEveryVideoCycle bool

	// console interface
	console console.UserInterface

	// channel for communicating with the debugger from the ctrl-c goroutine
	intChan chan os.Signal

	// channel for communicating with the debugger from the gui goroutine
	guiChan chan gui.Event

	// record user input to a script file
	scriptScribe script.Scribe

	// \/\/\/ inputLoop \/\/\/

	// buffer for user input
	input []byte

	// any error from previous emulation step
	lastStepError bool

	// we accumulate break, trap and watch messsages until we can service them
	// if the strings are empty then no break/trap/watch event has occurred
	breakMessages string
	trapMessages  string
	watchMessages string

	// halt the emulation (but not the debugger)
	haltEmulation bool

	// continue the emulation
	continueEmulation bool

	// \/\/\/ currently these run related booleans are set by various commands
	// and used by the input loop. however, I think it would be better if they
	// were return conditions from the parseInput() function

	// whether the debugger is to continue with the debugging loop
	// set to false only when debugger is to finish
	running bool

	// continue emulation until a halt condition is encountered
	runUntilHalt bool
}

// NewDebugger creates and initialises everything required for a new debugging
// session. Use the Start() method to actually begin the session.
func NewDebugger(tv television.Television, scr gui.GUI) (*Debugger, error) {
	var err error

	dbg := &Debugger{tv: tv, scr: scr}

	dbg.digest, err = screendigest.NewSHA1(dbg.tv)
	if err != nil {
		return nil, errors.New(errors.DebuggerError, err)
	}

	// create a new VCS instance
	dbg.vcs, err = hardware.NewVCS(dbg.tv)
	if err != nil {
		return nil, errors.New(errors.DebuggerError, err)
	}

	// create instance of disassembly -- the same base structure is used
	// for disassemblies subseuquent to the first one.
	dbg.disasm = &disassembly.Disassembly{}

	// set up debugging interface to memory
	dbg.dbgmem = &memoryDebug{mem: dbg.vcs.Mem, symtable: &dbg.disasm.Symtable}

	// set up reflection monitor
	dbg.relfectMonitor = reflection.NewMonitor(dbg.vcs, dbg.scr)
	dbg.relfectMonitor.Activate(true)

	// set up breakpoints/traps
	dbg.breakpoints = newBreakpoints(dbg)
	dbg.traps = newTraps(dbg)
	dbg.watches = newWatches(dbg)
	dbg.stepTraps = newTraps(dbg)

	// default ONHALT command sequence
	dbg.commandOnHaltStored = defaultOnHalt

	// default ONSTEP command sequnce
	dbg.commandOnStep = defaultOnStep
	dbg.commandOnStepStored = dbg.commandOnStep

	// make synchronisation channels
	dbg.intChan = make(chan os.Signal, 1)
	dbg.guiChan = make(chan gui.Event, 2)
	signal.Notify(dbg.intChan, os.Interrupt)

	// connect debugger to gui
	dbg.scr.SetEventChannel(dbg.guiChan)

	// allocate memory for user input
	dbg.input = make([]byte, 255)

	return dbg, nil
}

// Start the main debugger sequence. starting the debugger is a distinct
// operation to creating the debugger.
func (dbg *Debugger) Start(cons console.UserInterface, initScript string, cartload cartridgeloader.Loader) error {
	// prepare user interface
	if cons == nil {
		dbg.console = &console.PlainTerminal{}
	} else {
		dbg.console = cons
	}

	err := dbg.console.Initialise()
	if err != nil {
		return errors.New(errors.DebuggerError, err)
	}
	defer dbg.console.CleanUp()

	dbg.console.RegisterTabCompleter(commandline.NewTabCompletion(debuggerCommands))

	err = dbg.loadCartridge(cartload)
	if err != nil {
		return errors.New(errors.DebuggerError, err)
	}

	dbg.running = true

	// run initialisation script
	if initScript != "" {
		dbg.console.Silence(true)

		plb, err := script.StartPlayback(initScript)
		if err != nil {
			dbg.print(console.StyleError, "error running debugger initialisation script: %s\n", err)
		}

		err = dbg.inputLoop(plb, false)
		if err != nil {
			dbg.console.Silence(false)
			return errors.New(errors.DebuggerError, err)
		}

		dbg.console.Silence(false)
	}

	// prepare and run main input loop. inputLoop will not return until
	// debugging session is to be terminated
	err = dbg.inputLoop(dbg.console, false)
	if err != nil {
		return errors.New(errors.DebuggerError, err)
	}
	return nil
}

// loadCartridge makes sure that the cartridge loaded into vcs memory and the
// available disassembly/symbols are in sync.
//
// NEVER call vcs.AttachCartridge except through this function
//
// this is the glue that hold the cartridge and disassembly packages
// together
func (dbg *Debugger) loadCartridge(cartload cartridgeloader.Loader) error {
	err := setup.AttachCartridge(dbg.vcs, cartload)
	if err != nil {
		return err
	}

	symtable, err := symbols.ReadSymbolsFile(cartload.Filename)
	if err != nil {
		dbg.print(console.StyleError, "%s", err)
		// continuing because symtable is always valid even if err non-nil
	}

	err = dbg.disasm.FromMemory(dbg.vcs.Mem.Cart, symtable)
	if err != nil {
		return err
	}

	err = dbg.vcs.TV.Reset()
	if err != nil {
		return err
	}

	return nil
}

// parseInput splits the input into individual commands. each command is then
// passed to parseCommand for processing
//
// interactive argument should be true if  the input that has just come from
// the user (ie. via an interactive terminal). only interactive input will be
// added to a new script file.
//
// returns a boolean stating whether the emulation should continue with the
// next step
func (dbg *Debugger) parseInput(input string, interactive bool, auto bool) (bool, error) {
	var result parseCommandResult
	var err error
	var continueEmulation bool

	// ignore comments
	if strings.HasPrefix(input, "#") {
		return false, nil
	}

	// divide input if necessary
	commands := strings.Split(input, ";")
	for i := 0; i < len(commands); i++ {

		// try to record command now if it is not a result of an "autocommand"
		// (ONSTEP, ONHALT). if there's an error as a result of parsing, it
		// will be rolled back before committing
		if !auto {
			dbg.scriptScribe.WriteInput(commands[i])
		}

		// parse command. format of command[i] wil be normalised
		result, err = dbg.parseCommand(&commands[i], interactive)
		if err != nil {
			dbg.scriptScribe.Rollback()
			return false, err
		}

		// the result from parseCommand() tells us what to do next
		switch result {
		case doNothing:
			// most commands don't require us to do anything
			break

		case stepContinue:
			// emulation should continue to next step
			continueEmulation = true

		case emptyInput:
			// input was empty. if this was an interactive input then try the
			// default step command
			if interactive {
				return dbg.parseInput(onEmptyInput, interactive, auto)
			}
			return false, nil

		case scriptRecordStarted:
			// command has caused input script recording to begin. rollback the
			// call to recordCommand() above because we don't want to record
			// the fact that we've starting recording in the script itsel
			dbg.scriptScribe.Rollback()

		case scriptRecordEnded:
			// nothing special required when script recording has completed
		}

	}

	return continueEmulation, nil
}
