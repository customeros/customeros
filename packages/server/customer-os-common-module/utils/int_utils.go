package utils

import (
	"math/rand"
	"time"
)

// Create a package-level random source
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomInt(min, max int) int {
	// Generate a random integer between min and max
	return rng.Intn(max-min+1) + min
}
