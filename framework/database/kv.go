package db

import (
	"fmt"

	kvstore "github.com/multiverse-os/starshipyard/framework/database/kvstore"
)

type KV struct {
	Store       *kvstore.Database
	Collections map[string]*kvstore.KeyValue
}

// TODO: should move router definition to a route.go file
func InitKV(path string) *KV {
	store, err := kvstore.New(path)
	if err != nil {
		panic(fmt.Sprintf("[fatal error] failed to open session DB: %v", err))
	}
	return &KV{
		Store:       store,
		Collections: map[string]*kvstore.KeyValue{},
	}
}

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
