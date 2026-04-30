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

func Random(len int) (int, int) {
	max = big.NewInt(100)
	chance, _ = rand.Int(rand.Reader, max)
	val := chance.Int64()

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

	one = len / parts

	return parts, one
}

func Parts(m int64) int {
	part, _ := rand.Int(rand.Reader, big.NewInt(m))
	return int(part.Int64())
}
