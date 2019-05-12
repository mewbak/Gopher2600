package debugger

import (
	"gopher2600/debugger/console"
	"strings"
)

// wrapper function for UserPrint(). useful for normalising the input string
// before passing to the real UserPrint. it also allows us to easily obey
// directives such as the silent directive without passing the burden onto UI
// implementors
func (dbg *Debugger) print(pp console.PrintProfile, s string, a ...interface{}) {
	// trim *all* trailing newlines - UserPrint() will add newlines if required
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return
	}

	dbg.console.UserPrint(pp, s, a...)

	// output to script file
	if pp.IncludeInScriptOutput() {
		dbg.scriptScribe.WriteOutput(s, a...)
	}
}

// convenient but inflexible alternative to print()
func (dbg *Debugger) Write(p []byte) (n int, err error) {
	dbg.print(console.Feedback, string(p))
	return len(p), nil
}
