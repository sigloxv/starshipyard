package datastore

import (
	"fmt"

	kvstore "github.com/multiverse-os/starshipyard/framework/datastore/kvstore"
)

// TODO: Build in a cache into the KV object, logging for rollback,
type KV struct {
	Store       *kvstore.Database
	Collections map[string]*kvstore.KeyValue
}

func OpenKVStore(path string) Datastore {
	store, err := kvstore.New(path)
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to open session datastore: %v", err))
	}
	return &KV{
		Store:       store,
		Collections: map[string]*kvstore.KeyValue{},
	}
}

func (self *KV) Close() {
	self.Close()
}

//=============================================================================

// NOTE: Chainable collection initialization for an easier to use API for simple
// uses of the KV store; see sessions
func (self *KV) WithCollection(name string) *KV {
	self.NewCollection(name)
	return self
}

// TODO: Or perhaps LoadOrCreateCollection(name string)
func (self *KV) NewCollection(name string) (*kvstore.KeyValue, error) {
	collection, err := self.Store.NewCollection(name)
	if err != nil {
		return nil, fmt.Errorf("[error] failed to create '", name, "' collection: ", err)
	}
	self.Collections[name] = collection
	return collection, nil
}
