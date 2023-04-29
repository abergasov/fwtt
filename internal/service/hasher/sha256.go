package hasher

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

func (s *Service) CheckSHA256(difficulty, nonce uint32, challenge, verifyHash string) bool {
	verifyBytes := s.verifyHash(difficulty, verifyHash)
	if verifyBytes == nil {
		return false
	}

	return bytes.Equal(s.EncodeSHA256(nonce, challenge), verifyBytes)
}

func (s *Service) EncodeSHA256(nonce uint32, challenge string) []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", challenge, nonce)))
	return hash[:]
}
