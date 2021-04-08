package record

import "github.com/multiverse-os/starshipyard/framework/datastore/database/id"

type Record struct {
	Collection *Collection
	Name       string
	Id         id.Id
	Key        Key
	Value      Value
}
