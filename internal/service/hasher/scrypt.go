package hasher

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/scrypt"
)

type ScryptConfig struct {
	N      int `json:"n"`
	R      int `json:"r"`
	P      int `json:"p"`
	KeyLen int `json:"key_len"`
}

func DefaultScryptConfig() ScryptConfig {
	return ScryptConfig{N: 1024, R: 1, P: 1, KeyLen: 32}
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
