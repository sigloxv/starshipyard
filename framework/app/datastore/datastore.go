package datastore

type OptionType uint8

const (
	Find OptionType = iota
	Where
	Update
	Create
	Save
	Destroy
)

type Option struct {
	Type          OptionType
	Modifications []*Modification
}

func (self *Option) And(options ...[]*Optiuon)

type Key []byte

func (self Key) String() string { return string(self) }

// Store as BSON in string or bytes?
// Then that gets converted to the struct?
type Value []byte

func (self Value) String() string { return string(self) }

type Store interface {
	// TODO: Review ActiveRecord modern documentation and implement functions
	// using the same name where they are good. The remaining functionality will
	// be custom; but this makes the transition from Rails much easier.

	// TODO: These should be returning an array of errors.

	// Database Functions
	Close() error

	// Collection Functions
	NewCollection(collectionName string) (collection *Collection, err error)
	DestroyCollection(collectionName string) (success bool, err error)

	DestroyAll(collectionName string) (success bool, err error)
	All(collectionName string) (records []*Record, err error)

	// Record Functions
	Find(options ...[]*Option) (record *Record, err error)
	Where(options ...[]*Option) (records []*Record, err error)

	Create(options ...[]*Option) (r *Record, saved bool, err error)
	Save(record *Record, options ...[]*Option) (record *Record, saved bool, err error)
	Update(record *Record, options ...[]*Option) (record *Record, saved bool, err error)
	// TODO: These can be methods of Record too for simplicity.
	Destroy(record *Record)
}

type DatabaseType uint8

const (
	KeyValueStore DatabaseType = iota
	DocumentStore
	GraphStore
	ColumnStore
	RowStore
)

type DatabasePermissions uint8

const (
	WriteOnly DatabasePermissions = iota
	ReadOnly
	ReadAndWrite
	Disabled
)

type Database struct {
	Type        DatabaseType
	Collections []*database.Collection

	Permissions database.Permissions

	WriteDatabase  *database.Store
	CacheDatabases []*database.Store

	History []*database.Action
}

type KeyValue struct {
}

type Document struct {
}

type Graph struct {
}

type Cache struct {
}

func OpenKV(path string) (*KV, error) {
	return nil, nil
}

func (self *KV) Close() error {
	return nil
}
