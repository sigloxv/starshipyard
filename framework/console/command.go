package console

import (
	"bytes"
	"fmt"
	"sort"
	"text/tabwriter"
)

type Cmd struct {
	Name     string
	Aliases  []string
	Func     func(self *Context)
	Help     string
	LongHelp string

	Completer func(args []string) []string
	children  map[string]*Cmd
}

func (self *Cmd) AddCmd(cmd *Cmd) {
	if self.children == nil {
		self.children = make(map[string]*Cmd)
	}
	self.children[cmd.Name] = cmd
}

func (self *Cmd) DeleteCmd(name string) {
	delete(self.children, name)
}

func (self *Cmd) Children() []*Cmd {
	var cmds []*Cmd
	for _, cmd := range self.children {
		cmds = append(cmds, cmd)
	}
	sort.Sort(cmdSorter(cmds))
	return cmds
}

func (self *Cmd) hasSubcommand() bool {
	if len(self.children) > 1 {
		return true
	}
	if _, ok := self.children["help"]; !ok {
		return len(self.children) > 0
	}
	return false
}

func (self Cmd) HelpText() string {
	var b bytes.Buffer
	p := func(s ...interface{}) {
		fmt.Fprintln(&b)
		if len(s) > 0 {
			fmt.Fprintln(&b, s...)
		}
	}
	if self.LongHelp != "" {
		p(self.LongHelp)
	} else if self.Help != "" {
		p(self.Help)
	} else if self.Name != "" {
		p(self.Name, "has no help")
	}
	if self.hasSubcommand() {
		p("Commands:")
		w := tabwriter.NewWriter(&b, 0, 4, 2, ' ', 0)
		for _, child := range self.Children() {
			fmt.Fprintf(w, "\t%s\t\t\t%s\n", child.Name, child.Help)
		}
		w.Flush()
		p()
	}
	return b.String()
}

func (self *Cmd) findChildCmd(name string) *Cmd {
	if cmd, ok := self.children[name]; ok {
		return cmd
	}
	// find alias matching the name
	for _, cmd := range self.children {
		for _, alias := range cmd.Aliases {
			if alias == name {
				return cmd
			}
		}
	}
	return nil
}

func (self Cmd) FindCmd(args []string) (*Cmd, []string) {
	var cmd *Cmd
	for i, arg := range args {
		if cmd1 := self.findChildCmd(arg); cmd1 != nil {
			cmd = cmd1
			self = *cmd
			continue
		}
		return cmd, args[i:]
	}
	return cmd, nil
}

type cmdSorter []*Cmd

func (self cmdSorter) Len() int           { return len(self) }
func (self cmdSorter) Less(i, j int) bool { return self[i].Name < self[j].Name }
func (self cmdSorter) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
