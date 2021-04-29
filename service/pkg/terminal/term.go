package terminal

import (
	"io"

	"golang.org/x/term"
)

func New(c io.ReadWriter) *term.Terminal {
	return term.NewTerminal(c, "> ")
}
