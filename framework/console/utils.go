package console

import (
	"github.com/abiosoft/readline"
)

func clearScreen(s *Shell) error {
	_, err := readline.ClearScreen(s.writer)
	return err
}
