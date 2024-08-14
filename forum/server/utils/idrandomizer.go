package util

import (
	"math/rand"
	"time"
)

// IDRandomizer gets a randomized string and index for the database.
//
// Usage:
//
// Pass the desired length of the ID as the only argument.
//
// user_ID, index := IDRandomizer(10) -> "elperjlnCr", 141
func IDRandomizer(length int) (string, int) {
	rand.Seed(time.Now().Unix())
	randomString := make([]byte, length)

	for i := 0; i < length; i++ {
		selection := rand.Intn(3-0) + 0
		switch selection {
		case 0:
			randomString = append(randomString, byte(97+rand.Intn(25)))
		case 1:
			randomString = append(randomString, byte(65+rand.Intn(25)))
		case 2:
			randomString = append(randomString, byte(48+rand.Intn(9)))
		}

	}
	return string(randomString), ModuloID(randomString)
}

// Helper function for IDRandomizer that converts the ID to Index
//
// Example:
// "elperjlnCr" -> 141
func ModuloID(id []byte) int {
	var index int
	var listSize int = 1000
	for _, bt := range id {
		index += int(bt)
	}
	index *= 11
	index %= listSize
	return index
}
