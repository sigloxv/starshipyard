package datastore

type DatastoreType int

const (
	KVStore DatastoreType = iota
	ObjectStore
	//ConstantStore
	//Timeseries
	//Trie
	//AppendOnly
)

// Aliases
const (
	DocumentStore = ObjectStore
)

type Datastore interface {
	Close()
}
