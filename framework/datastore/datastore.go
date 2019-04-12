package datastore

type DatastoreType int

const (
	KVStore DatastoreType = iota
	Session
	Cache
	ObjectStore
	//Graph
	//Blob
	//FileStore
	//ConstantStore
	//Timeseries
	//Trie
	//AppendOnly
	//Columnar (should have SQL aliased)
)

// Aliases
const (
	DocumentStore = ObjectStore
)

type Datastore interface {
	Close()
}
