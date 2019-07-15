package debugger

import (
	"fmt"
	"gopher2600/debugger/console"
	"strings"
)

// wrapper function for UserPrint(). useful for normalising the input string
// before passing to the real UserPrint. it also allows us to easily obey
// directives such as the silent directive without passing the burden onto UI
// implementors
func (dbg *Debugger) print(sty console.Style, s string, a ...interface{}) {
	// resolve string placeholders and return if the resulting string is empty
	if sty != console.StyleHelp {
		s = fmt.Sprintf(s, a...)
		if len(s) == 0 {
			return
		}
	}

	// trim *all* trailing newlines - UserPrint() will add newlines if required
	s = strings.TrimRight(s, "\n")
	if len(s) == 0 {
		return
	}

	dbg.console.UserPrint(sty, s)

	// output to script file
	if sty.IncludeInScriptOutput() {
		dbg.scriptScribe.WriteOutput(s)
	}
}

// styleWriter is a wrapper for Debugger.print(). the result of
// printStyle() can be used as an implementation of the io.Writer interface
type styleWriter struct {
	dbg   *Debugger
	style console.Style
}

func (dbg *Debugger) printStyle(sty console.Style) *styleWriter {
	return &styleWriter{
		dbg:   dbg,
		style: sty,
	}
}

// convenient but inflexible alternative to print()
func (wrt styleWriter) Write(p []byte) (n int, err error) {
	wrt.dbg.print(wrt.style, string(p))
	return len(p), nil
}
