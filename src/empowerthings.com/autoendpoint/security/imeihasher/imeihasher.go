package imeihasher

import (
	"errors"
	//"crypto"
	"crypto/sha256"
	encoding "encoding/base64"
	//"empowerthings.com/cumulis/utils/encoding"
	//"fmt"
)

// DEFAULT_LOOPS are the times we are hashing the string
const DEFAULT_LOOPS = 1000

// ImeiHasher is the struct that holds the imeihasher data
type ImeiHasher struct {
	secret []byte
	loops  int
}

// NewImeiHasher returns a new ImeiHasher instance with DEFAULT_LOOPS
func NewImeiHasher(secret string) (*ImeiHasher, error) {
	if len(secret) == 0 {
		return nil, errors.New("Missing secret string!")
	}
	byteSecret := []byte(secret)
	return &ImeiHasher{secret: byteSecret, loops: DEFAULT_LOOPS}, nil
}

// NewImeiHasherWithLoops returns a new ImeiHasher specifying the loops to hash
func NewImeiHasherWithLoops(secret string, loops int) (*ImeiHasher, error) {
	if len(secret) == 0 {
		return nil, errors.New("Missing secret string!")
	}
	if loops == 0 {
		loops = 1
	}
	byteSecret := []byte(secret)
	return &ImeiHasher{secret: byteSecret, loops: loops}, nil
}

// Hash performs the hash operation
func (ih *ImeiHasher) Hash(imei string) string {
	hashedImei := []byte(imei)
	for i := 0; i < ih.loops; i++ {
		hashedImei = []byte(string(hashedImei) + string(ih.secret))
		h := sha256.Sum256(hashedImei)
		hashedImei = h[:]
	}
	return encoding.RawURLEncoding.EncodeToString(hashedImei)
}
