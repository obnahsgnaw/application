package security

import "crypto"

// Hash cal hash
func Hash(data []byte, hs crypto.Hash) (hashed []byte, err error) {
	hash := hs.New()

	if _, err = hash.Write(data); err != nil {
		return
	}

	return hash.Sum(nil), nil
}
