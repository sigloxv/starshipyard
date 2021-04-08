package job

import (
	"time"
)

type StatusType int

const (
	Completed StatusType = iota
	Failed
	Running
	Waiting
)

type Context struct {
}

type Action func(context *Context) error

type Job struct {
	Name string

	Status StatusType

	Unique bool
	Repeat bool

	RunAt time.Time

	Action *Action
}
