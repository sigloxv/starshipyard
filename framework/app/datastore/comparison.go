package datastore

import data "github.com/multiverse-os/starshipyard/framework/datastore/datatypes"

type CompareType uint8

const (
	Is CompareType = iota
	IsNot
	IsLess
	IsLessOrEqual
	IsMore
	IsMoreOrEqual
)

// Aliasing CompareType
const (
	IsEqual          = Is
	IsNotEqual       = IsNot
	IsGreater        = IsMore
	IsGreaterOrEqual = IsMoreOrEqual
)

type Comparison struct {
	Field object.StructField
	Type  CompareType
	Value data.Value
}
