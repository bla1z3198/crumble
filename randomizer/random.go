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
	last   int
)

func Random(len int) (int, int, int) {
	max = big.NewInt(100)
	chance, _ = rand.Int(rand.Reader, max)

	switch {
	case chance.Int64() <= 70 && chance.Int64() >= 22:
		parts = 10
	case chance.Int64() <= 22 && chance.Int64() >= 6:
		parts = 56
	case chance.Int64() <= 6 && chance.Int64() >= 4:
		parts = 2
	case chance.Int64() <= 4 && chance.Int64() >= 0:
		parts = 4
	default:
		parts = 20
	}

	one = len / parts
	last = len - one

	return parts, one, last
}
