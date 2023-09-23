package graph

import (
	"math/rand"
	"time"
)

var (
	charset = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
)

func generateNewRandomCustomerOsId() string {
	customerOsID := "C-" + generateRandomStringFromCharset(3) + "-" + generateRandomStringFromCharset(3)
	return customerOsID
}

func generateRandomStringFromCharset(length int) string {
	// Create a new source based on the current time's Unix timestamp (in nanoseconds)
	source := rand.NewSource(time.Now().UnixNano())
	// Initialize a random number generator (RNG) with the source
	rng := rand.New(source)

	var output string
	for i := 0; i < length; i++ {
		randChar := charset[rng.Intn(len(charset))]
		output += string(randChar)
	}
	return output
}
