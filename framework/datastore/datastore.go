package datastore

// TODO: Datastore open funcitons should return the close function hneeded to
// properly close the database. This can then be stored in the shutdown slice of
// functions for clean shutdown

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
