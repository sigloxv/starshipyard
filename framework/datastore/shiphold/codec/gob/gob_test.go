package gob

import (
	"testing"

	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/codec/cursor"
)

func TestGob(t *testing.T) {
	cursor.RoundtripTester(t, Codec)
}
