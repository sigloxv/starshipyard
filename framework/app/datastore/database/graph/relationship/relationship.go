package relationship

import (
	"github.com/multiverse-os/starshipyard/framework/datastore/database/id"
)

type Relationship struct {
	EntityId   id.Id
	Components []*Component
}

// TODO: Because we can not easilyn access entity, this may not be a good choice
//       for the design. We may need to break off entity or merge this back into
//       graph for simplicity.
