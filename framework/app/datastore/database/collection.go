package database

import (
	"time"
)

type SortDirection uint8

const (
	Ascending SortDirection = iota
	Descending
)

type Sort struct {
	Column    string
	Direction SortDirection
}

type Index struct {
	Name     string
	Sorting  []*Sort
	IsSorted bool
	Records  []*Record
}

type Collection struct {
	Name          string
	Records       []*Record
	Indexes       []*Index
	LastUpdatedAt time.Time
	History       chan *Record
}

func (self *Collection) StepBackwards() (success bool, err error) {
	// TODO: Take last item from the history, remove(pop) it. Then use this item
	//       to revert the last change.

	// TODO: Record should also have all its own updates stored within it, in
	//       order so that it can be used to revert changes or show historical
	//       versions using BSON diffing.
}

func (self *Collection) All() (records []*Record, err error) {
	return self.Records
}
