package randomizer

import (
	"crypto/rand"
	"math/big"
)

var (
	parts  int
	one    int
	max    *big.Int
	chance *big.Int
)

func Random(len int, ch chan []int) {
	// Max value
	max = big.NewInt(100)
	// Rand int
	chance, _ = rand.Int(rand.Reader, max)
	// Convert to int64
	val := chance.Int64()
	// Compare val with fixed values
	switch {
	case val > 30:
		parts = 7 + Parts(6)
	case val > 80:
		parts = 3 + Parts(4)
	case val > 92:
		parts = 15 + Parts(11)
	default:
		parts = 10 + Parts(5)
	}
	// Define size of one part
	one = len / parts
	// Write results
	result := make([]int, 2)
	result[0] = parts
	result[1] = one
	// Send results
	ch <- result
}

func Parts(m int64) int {
	part, _ := rand.Int(rand.Reader, big.NewInt(m))
	return int(part.Int64())
}
