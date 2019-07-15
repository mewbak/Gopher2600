package console

import (
	"fmt"
	"gopher2600/gui"
	"io"
	"os"
)

// PlainTerminal is the default, most basic terminal interface
type PlainTerminal struct {
	input  io.Reader
	output io.Writer
}

// Initialise perfoms any setting up required for the terminal
func (pt *PlainTerminal) Initialise() error {
	pt.input = os.Stdin
	pt.output = os.Stdout
	return nil
}

// CleanUp perfoms any cleaning up required for the terminal
func (pt *PlainTerminal) CleanUp() {
}

// RegisterTabCompleter adds an implementation of TabCompleter to the terminal
func (pt *PlainTerminal) RegisterTabCompleter(TabCompleter) {
}

// UserPrint is the plain terminal print routine
func (pt PlainTerminal) UserPrint(pp Style, s string, a ...interface{}) {
	switch pp {
	case StyleError:
		s = fmt.Sprintf("* %s", s)
	case StyleHelp:
		s = fmt.Sprintf("  %s", s)
	}

	s = fmt.Sprintf(s, a...)
	pt.output.Write([]byte(s))

	if pp != StylePrompt {
		pt.output.Write([]byte("\n"))
	}
}

// UserRead is the plain terminal read routine
func (pt PlainTerminal) UserRead(input []byte, prompt Prompt, _ chan gui.Event, _ func(gui.Event) error) (int, error) {
	pt.UserPrint(prompt.Style, prompt.Content)

	n, err := pt.input.Read(input)
	if err != nil {
		return n, err
	}
	return n, nil
}

// IsInteractive satisfies the console.UserInput interface
func (pt *PlainTerminal) IsInteractive() bool {
	return true
}
