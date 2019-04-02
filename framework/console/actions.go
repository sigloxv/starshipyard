package console

import (
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

type Actions interface {
	ReadLine() string
	ReadLineErr() (string, error)
	ReadLineWithDefault(string) string
	ReadPassword() string
	ReadPasswordErr() (string, error)
	ReadMultiLinesFunc(f func(string) bool) string
	ReadMultiLines(terminator string) string
	Println(val ...interface{})
	Print(val ...interface{})
	Printf(format string, val ...interface{})
	ShowPaged(text string) error
	ShowPagedReader(r io.Reader) error
	MultiChoice(options []string, text string) int
	Checklist(options []string, text string, init []int) []int
	SetPrompt(prompt string)
	SetMultiPrompt(prompt string)
	SetMultiChoicePrompt(prompt, spacer string)
	SetChecklistOptions(open, selected string)
	ShowPrompt(show bool)
	Cmds() []*Cmd
	HelpText() string
	ClearScreen() error
	Stop()
}

type consoleActionsImpl struct {
	*Console
}

// ReadLine reads a line from standard input.
func (self *consoleActionsImpl) ReadLine() string {
	line, _ := self.readLine()
	return line
}

func (self *consoleActionsImpl) ReadLineErr() (string, error) {
	return self.readLine()
}

func (self *consoleActionsImpl) ReadLineWithDefault(defaultValue string) string {
	self.reader.defaultInput = defaultValue
	line, _ := self.readLine()
	self.reader.defaultInput = ""
	return line
}

func (self *consoleActionsImpl) ReadPassword() string {
	return self.reader.readPassword()
}

func (self *consoleActionsImpl) ReadPasswordErr() (string, error) {
	return self.reader.readPasswordErr()
}

func (self *consoleActionsImpl) ReadMultiLinesFunc(f func(string) bool) string {
	lines, _ := self.readMultiLinesFunc(f)
	return lines
}

func (self *consoleActionsImpl) ReadMultiLines(terminator string) string {
	return self.ReadMultiLinesFunc(func(line string) bool {
		if strings.HasSuffix(strings.TrimSpace(line), terminator) {
			return false
		}
		return true
	})
}

func (self *consoleActionsImpl) Println(val ...interface{}) {
	self.reader.buf.Truncate(0)
	fmt.Fprintln(self.writer, val...)
}

func (self *consoleActionsImpl) Print(val ...interface{}) {
	self.reader.buf.Truncate(0)
	fmt.Fprint(self.reader.buf, val...)
	fmt.Fprint(self.writer, val...)
}

func (self *consoleActionsImpl) Printf(format string, val ...interface{}) {
	self.reader.buf.Truncate(0)
	fmt.Fprintf(self.reader.buf, format, val...)
	fmt.Fprintf(self.writer, format, val...)
}

func (self *consoleActionsImpl) MultiChoice(options []string, text string) int {
	choice := self.multiChoice(options, text, nil, false)
	return choice[0]
}
func (self *consoleActionsImpl) Checklist(options []string, text string, init []int) []int {
	return self.multiChoice(options, text, init, true)
}
func (self *consoleActionsImpl) SetPrompt(prompt string) {
	self.reader.prompt = prompt
	self.reader.scanner.SetPrompt(self.reader.rlPrompt())
}

func (self *consoleActionsImpl) SetMultiPrompt(prompt string) {
	self.reader.multiPrompt = prompt
}

func (self *consoleActionsImpl) SetMultiChoicePrompt(prompt, spacer string) {
	strMultiChoice = prompt
	strMultiChoiceSpacer = spacer
}
func (self *consoleActionsImpl) SetChecklistOptions(open, selected string) {
	strMultiChoiceOpen = open
	strMultiChoiceSelect = selected
}

func (self *consoleActionsImpl) ShowPrompt(show bool) {
	self.reader.showPrompt = show
	self.reader.scanner.SetPrompt(self.reader.rlPrompt())
}

func (self *consoleActionsImpl) Cmds() []*Cmd {
	var cmds []*Cmd
	for _, cmd := range self.rootCmd.children {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (self *consoleActionsImpl) ClearScreen() error {
	return clearScreen(self.Console)
}

func (self *consoleActionsImpl) ShowPaged(text string) error {
	return showPagedReader(self.Console, strings.NewReader(text))
}

func (self *consoleActionsImpl) ShowPagedReader(r io.Reader) error {
	return showPagedReader(self.Console, r)
}

func (self *consoleActionsImpl) Stop() {
	self.stop()
}

func (self *consoleActionsImpl) HelpText() string {
	return self.rootCmd.HelpText()
}

func showPagedReader(self *Console, r io.Reader) error {
	var cmd *exec.Cmd
	if self.pager == "" {
		self.pager = "less"
	}
	cmd = exec.Command(self.pager, self.pagerArgs...)
	cmd.Stdout = self.writer
	cmd.Stderr = self.writer
	cmd.Stdin = r
	return cmd.Run()
}
