package sessions

import (
	store "github.com/multiverse-os/starshipyard/framework/sessions/store"
)

type Store struct {
	Count    int
	Sessions *store.DB
	Encoder  func(string) []byte
	Decoder  func([]byte) string
}

func NewStore(path string) *Store {
	sessions, err := store.Open("sessions.store", nil)
	if err != nil {
		panic(err)
	}
	return &Store{
		Sessions: sessions,
	}
}
