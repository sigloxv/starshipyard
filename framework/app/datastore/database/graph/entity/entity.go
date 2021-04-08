package entity

import (
	id "github.com/multiverse-os/starshipyard/framework/datastore/database/id"
	relationship "github.com/multiverse-os/starshipyard/framework/datastore/types/relationship"
)

type Entity struct {
	Id id.Id

	Record     *database.Record
	Collection *database.Collection
	Database   *database.Database

	Relationships []*Relationship
}

// NOTE: As in Rails, new creates the object without saving,
//       while create will create the object and save it.

func New(collection *database.Collection, record *database.Reord, relationships ...relationship.Relationship) (entity *Entity, err error) {
	// TODO: Look up the record or error out (no creating new ones), then save the
	//       relationship data for this reocrd inside an entity. The entities will
	//       be stored in a key value storage.
}
