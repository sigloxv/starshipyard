package model

import "time"

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

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Model struct {
	*Hooks // Embedded so they are called directly on Model
	Timestamps

	ID   string // Preferably use BSON
	Name string

	Object interface{}
}

type Query struct {
	Attribute string
	Value     string
}

// TODO: If we use just a KV database, we can flatten the BSON and insert it by
// ID. Then store slices where the ids are sorted by createdAt, etc.
type ModelInterface interface {
	// TODO: Add ability to convert to JSON, BSON, etc.
	JSON() string
	BSON() string

	Where(query ...string)

	All() []*Model

	Paginate(page, perPage int) []*Model
	Search(search string) []*Model
	Random(records int) []*Model

	First() *Model
	Last() *Model
}
