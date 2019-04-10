package json

import (
	"testing"

	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/cursor"
)

func TestJSON(t *testing.T) {
	cursor.RoundtripTester(t, Codec)
}
