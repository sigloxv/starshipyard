package session

import (
	"time"

	id "github.com/multiverse-os/starshipyard/framework/datastore/id"
)

type Session struct {
	Id       id.Id
	ExpireAt time.Time
}

func New(expireAt time.Time) *Session {
	return &Session{
		ExpireAt: expireAt,
	}
}
