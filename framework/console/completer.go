package console

import (
	"strings"

	sh "github.com/multiverse-os/starshipyard/framework/console/sh"
)

type Completer struct {
	cmd      *Cmd
	disabled func() bool
}

func (self Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if self.disabled != nil && self.disabled() {
		return nil, len(line)
	}
	var words []string
	if w, err := sh.Split(string(line)); err == nil {
		words = w
	} else {
		// fall back
		words = strings.Fields(string(line))
	}

	var cWords []string
	prefix := ""
	if len(words) > 0 && pos > 0 && line[pos-1] != ' ' {
		prefix = words[len(words)-1]
		cWords = self.getWords(words[:len(words)-1])
	} else {
		cWords = self.getWords(words)
	}

	var suggestions [][]rune
	for _, w := range cWords {
		if strings.HasPrefix(w, prefix) {
			suggestions = append(suggestions, []rune(strings.TrimPrefix(w, prefix)))
		}
	}
	if len(suggestions) == 1 && prefix != "" && string(suggestions[0]) == "" {
		suggestions = [][]rune{[]rune(" ")}
	}
	return suggestions, len(prefix)
}

func (self Completer) getWords(w []string) (s []string) {
	cmd, args := self.cmd.FindCmd(w)
	if cmd == nil {
		cmd, args = self.cmd, w
	}
	if cmd.Completer != nil {
		return cmd.Completer(args)
	}
	for k := range cmd.children {
		s = append(s, k)
	}
	return
}
