package hasher

import (
	"encoding/hex"
)

type Service struct {
	scryptConfig ScryptConfig
}

func NewService() *Service {
	return &Service{
		scryptConfig: DefaultScryptConfig(),
	}
}

func (s *Service) ScryptConfig() *ScryptConfig {
	return &s.scryptConfig
}

func (s *Service) verifyHash(difficulty uint32, verifyHash string) []byte {
	if difficulty < 1 {
		difficulty = 1
	}
	if len(verifyHash) < int(difficulty) {
		return nil
	}
	// check that first difficulty symbols are zeros
	for i := 0; i < int(difficulty); i++ {
		if verifyHash[i] != '0' {
			return nil
		}
	}
	verifyBytes, err := hex.DecodeString(verifyHash)
	if err != nil {
		return nil
	}
	return verifyBytes
}
