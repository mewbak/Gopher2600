package main

import (
	"flag"
	"fmt"
	"gopher2600/debugger"
	"gopher2600/debugger/colorterm"
	"gopher2600/debugger/console"
	"gopher2600/disassembly"
	"gopher2600/errors"
	"gopher2600/performance"
	"gopher2600/playmode"
	"gopher2600/recorder"
	"gopher2600/regression"
	"io"
	"os"
	"path"
	"strings"
)

const defaultInitScript = ".gopher2600/debuggerInit"

func main() {
	progName := path.Base(os.Args[0])

	var mode string
	var modeArgPos int
	var modeFlags *flag.FlagSet
	var modeFlagsParse func()

	progModes := []string{"RUN", "PLAY", "DEBUG", "DISASM", "PERFORMANCE", "REGRESS"}
	defaultMode := "RUN"

	progFlags := flag.NewFlagSet(progName, flag.ContinueOnError)

	// prevent Parse() (part of the flag package) from outputting it's own error messages
	progFlags.SetOutput(&nopWriter{})

	err := progFlags.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			fmt.Printf("available modes: %s\n", strings.Join(progModes, ", "))
			fmt.Printf("default: %s\n", defaultMode)
			os.Exit(2)
		}

		// flags have been set that are not recognised. default to the RUN mode
		// and try again
		mode = defaultMode
		modeArgPos = 0
		modeFlags = flag.NewFlagSet(fmt.Sprintf("%s %s", progName, mode), flag.ExitOnError)
		modeFlagsParse = func() {
			if len(progFlags.Args()) >= modeArgPos {
				modeFlags.Parse(os.Args[1:])
			}
		}
	} else {
		switch progFlags.NArg() {
		case 0:
			// no arguments at all. suggest that a cartridge is required
			fmt.Println("* 2600 cartridge required")
			os.Exit(2)
		case 1:
			// a single argument has been supplied. assume it's a cartridge
			// name and set the mode to the default mode ...
			mode = defaultMode
			modeArgPos = 0

			// ... unless it apears in the list of modes. in which case, the
			// single argument is a specified mode. let the mode switch below
			// handle what to do next.
			arg := strings.ToUpper(progFlags.Arg(0))
			for i := range progModes {
				if progModes[i] == arg {
					mode = arg
					modeArgPos = 1
					break
				}
			}
		default:
			// many arguments have been supplied. the first argument must be
			// the mode (the switch below will compalin if it's invalid)
			mode = strings.ToUpper(progFlags.Arg(0))
			modeArgPos = 1
		}

		// all modes can have their own sets of flags. the following prepares
		// the foundations.
		modeFlags = flag.NewFlagSet(fmt.Sprintf("%s %s", progName, mode), flag.ExitOnError)
		modeFlagsParse = func() {
			if len(progFlags.Args()) >= modeArgPos {
				modeFlags.Parse(progFlags.Args()[modeArgPos:])
			}
		}
	}

	switch mode {
	default:
		fmt.Printf("* %s mode unrecognised\n", mode)
		os.Exit(2)

	case "RUN":
		fallthrough

	case "PLAY":
		tvType := modeFlags.String("tv", "AUTO", "television specification: NTSC, PAL")
		scaling := modeFlags.Float64("scale", 3.0, "television scaling")
		stable := modeFlags.Bool("stable", true, "wait for stable frame before opening display")
		record := modeFlags.Bool("record", false, "record user input to a file")
		recording := modeFlags.String("recording", "", "the file to use for recording/playback")
		modeFlagsParse()

		switch len(modeFlags.Args()) {
		case 0:
			if *recording == "" {
				fmt.Println("* 2600 cartridge required")
				os.Exit(2)
			}
			fallthrough
		case 1:
			err := playmode.Play(modeFlags.Arg(0), *tvType, float32(*scaling), *stable, *recording, *record)
			if err != nil {
				fmt.Printf("* %s\n", err)
				os.Exit(2)
			}
			if *record == true {
				fmt.Println("! recording completed")
			} else if *recording != "" {
				fmt.Println("! playback completed")
			}
		default:
			fmt.Printf("* too many arguments for %s mode\n", mode)
			os.Exit(2)
		}

	case "DEBUG":
		tvType := modeFlags.String("tv", "AUTO", "television specification: NTSC, PAL")
		termType := modeFlags.String("term", "COLOR", "terminal type to use in debug mode: COLOR, PLAIN")
		initScript := modeFlags.String("initscript", defaultInitScript, "terminal type to use in debug mode: COLOR, PLAIN")
		modeFlagsParse()

		dbg, err := debugger.NewDebugger(*tvType)
		if err != nil {
			fmt.Printf("* %s\n", err)
			os.Exit(2)
		}

		// start debugger with choice of interface and cartridge
		var term console.UserInterface

		switch strings.ToUpper(*termType) {
		default:
			fmt.Printf("! unknown terminal type (%s) defaulting to plain\n", *termType)
			fallthrough
		case "PLAIN":
			term = nil
		case "COLOR":
			term = new(colorterm.ColorTerminal)
		}

		switch len(modeFlags.Args()) {
		case 0:
			// it's okay if DEBUG mode is started with no cartridges
			fallthrough
		case 1:
			err := dbg.Start(term, *initScript, modeFlags.Arg(0))
			if err != nil {
				fmt.Printf("* %s\n", err)
				os.Exit(2)
			}
		default:
			fmt.Printf("* too many arguments for %s mode\n", mode)
			os.Exit(2)
		}

	case "DISASM":
		modeFlagsParse()

		switch len(modeFlags.Args()) {
		case 0:
			fmt.Println("* 2600 cartridge required")
			os.Exit(2)
		case 1:
			dsm, err := disassembly.FromCartrige(modeFlags.Arg(0))
			if err != nil {
				switch err.(type) {
				case errors.FormattedError:
					// print what disassembly output we do have
					if dsm != nil {
						dsm.Dump(os.Stdout)
					}
				}
				fmt.Printf("* %s\n", err)
				os.Exit(2)
			}
			dsm.Dump(os.Stdout)
		default:
			fmt.Printf("* too many arguments for %s mode\n", mode)
			os.Exit(2)
		}

	case "PERFORMANCE":
		display := modeFlags.Bool("display", false, "display TV output: boolean")
		tvType := modeFlags.String("tv", "AUTO", "television specification: NTSC, PAL")
		scaling := modeFlags.Float64("scale", 3.0, "television scaling")
		runTime := modeFlags.String("time", "5s", "run duration (note: there is a 2s overhead)")
		profile := modeFlags.Bool("profile", false, "perform cpu and memory profiling")
		modeFlagsParse()

		switch len(modeFlags.Args()) {
		case 0:
			fmt.Println("* 2600 cartridge required")
			os.Exit(2)
		case 1:
			err := performance.Check(os.Stdout, *profile, modeFlags.Arg(0), *display, *tvType, float32(*scaling), *runTime)
			if err != nil {
				fmt.Printf("* %s\n", err)
				os.Exit(2)
			}
		default:
			fmt.Printf("* too many arguments for %s mode\n", mode)
			os.Exit(2)
		}

	case "REGRESS":
		subMode := strings.ToUpper(progFlags.Arg(1))
		modeArgPos++
		switch subMode {
		default:
			modeArgPos-- // undo modeArgPos adjustment
			fallthrough

		case "RUN":
			// no additional arguments
			modeFlagsParse()
			err := regression.RegressRunTests(os.Stdout, modeFlags.Args())
			if err != nil {
				fmt.Printf("* error during regression tests: %s\n", err)
				os.Exit(2)
			}

		case "LIST":
			// no additional arguments
			modeFlagsParse()
			switch len(modeFlags.Args()) {
			case 0:
				err := regression.RegressList(os.Stdout)
				if err != nil {
					fmt.Printf("* error during regression listing: %s\n", err)
					os.Exit(2)
				}
			default:
				fmt.Printf("* no additional arguments required when using %s/%s mode\n", mode, subMode)
				os.Exit(2)
			}

		case "DELETE":
			answerYes := modeFlags.Bool("yes", false, "answer yes to confirmation")
			modeFlagsParse()

			switch len(modeFlags.Args()) {
			case 0:
				fmt.Println("* database key required (use REGRESS LIST to view)")
				os.Exit(2)
			case 1:

				// use stdin for confirmation unless "yes" flag has been sent
				var confirmation io.Reader
				if *answerYes == true {
					confirmation = new(yesReader)
				} else {
					confirmation = os.Stdin
				}

				err := regression.RegressDelete(os.Stdout, confirmation, modeFlags.Arg(0))
				if err != nil {
					fmt.Printf("* error deleting regression test: %s\n", err)
					os.Exit(2)
				}
			default:
				fmt.Printf("* only one entry can be deleted at at time when using %s/%s \n", mode, subMode)
				os.Exit(2)
			}

		case "ADD":
			tvType := modeFlags.String("tv", "AUTO", "television specification: NTSC, PAL (cartridge args only)")
			numFrames := modeFlags.Int("frames", 10, "number of frames to run (cartridge args only)")
			modeFlagsParse()

			switch len(modeFlags.Args()) {
			case 0:
				fmt.Println("* 2600 cartridge or playback file required")
				os.Exit(2)
			case 1:
				var rec regression.Handler

				if recorder.IsPlaybackFile(modeFlags.Arg(0)) {
					rec = &regression.PlaybackRegression{
						Script: modeFlags.Arg(0),
					}
				} else {
					rec = &regression.FrameRegression{
						CartFile:  modeFlags.Arg(0),
						TVtype:    *tvType,
						NumFrames: *numFrames}
				}

				err := regression.RegressAdd(os.Stdout, rec)
				if err != nil {
					fmt.Printf("* error adding regression test: %s\n", err)
					os.Exit(2)
				}
			default:
				fmt.Printf("* regression tests must be added one at a time when using %s/%s mode\n", mode, subMode)
				os.Exit(2)
			}
		}
	}
}

// special purpose io.Reader / io.Writer

type nopWriter struct{}

func (*nopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type yesReader struct{}

func (*yesReader) Read(p []byte) (n int, err error) {
	p[0] = 'y'
	return 1, nil
}
