package hasher

import (
	"bytes"
	"encoding/binary"
	"fwtt/internal/entites"

	"golang.org/x/crypto/scrypt"
)

func DefaultScryptConfig() entites.ScryptConfig {
	return entites.ScryptConfig{N: 1024, R: 1, P: 1, KeyLen: 32}
}

func (s *Service) CheckScrypt(difficulty, nonce uint32, challenge, verifyHash string) bool {
	verifyBytes := s.verifyHash(difficulty, verifyHash)
	if verifyBytes == nil {
		return false
	}
	return bytes.Equal(s.EncodeScrypt(nonce, challenge), verifyBytes)
}

func (s *Service) EncodeScrypt(nonce uint32, challenge string) []byte {
	salt := make([]byte, 4)
	binary.LittleEndian.PutUint32(salt, nonce)
	hash, err := scrypt.Key([]byte(challenge), salt, s.scryptConfig.N, s.scryptConfig.R, s.scryptConfig.P, s.scryptConfig.KeyLen)
	if err != nil {
		return []byte{}
	}
	return hash
}
