package database

import (
	"time"
)

type Action struct {
	Database   *Database
	Collection *Collection
	Record     *Record

	CreatedAt time.Time
	// TODO: Modification, only diff of what has been changed
}

type DatabaseType uint8

const (
	KV DatabaseType = iota
	Document
	ColumnOriented
	RowOriented
	Graph
)

type Database struct {
	IsCache     bool
	IsReadable  bool
	IsWriteable bool
	Type        DatabaseType

	Collections []*Collection // Collection should be table, bucket, or generic
	Records     uint64
	History     []*Action
}
