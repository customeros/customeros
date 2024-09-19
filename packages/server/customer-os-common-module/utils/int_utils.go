package utils

import (
	"math/rand"
	"time"
)

func GenerateRandomInt(min, max int) int {
	// Seed the random number generator to ensure different outputs each time
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between min and max
	return rand.Intn(max-min+1) + min
}
