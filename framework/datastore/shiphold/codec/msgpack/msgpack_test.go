package msgpack

import (
	"testing"

	cursor "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/cursor"
)

func TestMsgpack(t *testing.T) {
	cursor.RoundtripTester(t, Codec)
}
