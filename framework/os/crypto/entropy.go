package framework

import (
	"math/rand"
	"time"
)

func seedRandom() {
	rand.Seed(time.Now().UTC().UnixNano())
}
