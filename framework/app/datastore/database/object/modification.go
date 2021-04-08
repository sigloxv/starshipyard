package object

import (
	data "github.com/multiverse-os/starshipyard/framework/datastore/datatypes"
)

type ModificationType uint8

const (
	CreateType ModificationType = iota
	UpdateType
	DestroyType
	DestroyAllType
)

type Modification struct {
	FieldName string
	DataType  data.Type
	Value     interface{}
}
