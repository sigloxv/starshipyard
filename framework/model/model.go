package model

//import (
//	"fmt"
//)

type Context struct {
}

type Action func(context *Context) error

type Hooks struct {
	BeforeSave []Action
	AfterSave  []Action

	BeforeCreate []Action
	AfterCreate  []Action

	BeforeDelete []Action
	AfterDelete  []Action
}

type Model struct {
	Hooks *Hooks
	Name  string

	Object interface{}
}

type ModelInterface interface {
	All() []*Model

	Paginate(page, perPage int) []*Model
	Search(search string) []*Model
	Random(records int) []*Model

	First() *Model
	Last() *Model
}
