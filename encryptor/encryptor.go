package encryptor

import (
	"crypto/sha256"

	"golang.org/x/crypto/chacha20"
)

var (
	key   []byte
	nonce []byte
)

func Secret() ([32]byte, []byte) {
	hash := sha256.Sum256([]byte("dfghfghfgh568758!"))

	n_hash := hash[:]
	nonce = make([]byte, 12)
	for i := range nonce {
		nonce[i] = n_hash[31-i]
	}
	return hash, nonce
}

func Encrypt(wrapped []byte) []byte {
	a, b := Secret()
	cipherInstance, _ := chacha20.NewUnauthenticatedCipher(a[:], b)

	out := make([]byte, len(wrapped))
	cipherInstance.XORKeyStream(out, wrapped)
	return out
}
