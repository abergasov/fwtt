package validator_test

import (
	"encoding/hex"
	"fwtt/internal/entites"
	"fwtt/internal/logger"
	"fwtt/internal/service/hasher"
	"fwtt/internal/service/validator"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	difficulty          = 1
	maxAllowed          = 3
	challengeExpiration = 500 * time.Millisecond
)

func TestService_VerifyChallenge(t *testing.T) {
	serviceHasher := hasher.NewService()
	appLog, err := logger.NewAppLogger("test")
	require.NoError(t, err)
	serviceValidator := validator.NewService(appLog, difficulty, maxAllowed, challengeExpiration, serviceHasher, nil)

	t.Run("should not allow if requested more times than allowed", func(t *testing.T) {
		challengeParams := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParams.Difficulty, challengeParams.Challenges[0])
		for i := 0; i < maxAllowed; i++ {
			require.True(t, serviceValidator.VerifyChallenge(nonce, challengeParams.Challenges[0], challengeHash), "iteration %d", i)
		}
		require.False(t, serviceValidator.VerifyChallenge(nonce, challengeParams.Challenges[0], challengeHash))
	})
	t.Run("should not allow if nonce is wrong", func(t *testing.T) {
		challengeParams := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParams.Difficulty, challengeParams.Challenges[0])
		require.False(t, serviceValidator.VerifyChallenge(nonce+1, challengeParams.Challenges[0], challengeHash))
	})
	t.Run("should not allow if challenge is wrong", func(t *testing.T) {
		challengeParams := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParams.Difficulty, challengeParams.Challenges[0])
		require.False(t, serviceValidator.VerifyChallenge(nonce, challengeParams.Challenges[0]+"a", challengeHash))
	})
	t.Run("should not allow if hash is wrong", func(t *testing.T) {
		challengeParams := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParams.Difficulty, challengeParams.Challenges[0])
		require.False(t, serviceValidator.VerifyChallenge(nonce, challengeParams.Challenges[0], challengeHash+"a"))
	})
	t.Run("should not allow if time is expired", func(t *testing.T) {
		challengeParams := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParams.Difficulty, challengeParams.Challenges[0])
		time.Sleep(challengeExpiration + 1*time.Second)
		require.False(t, serviceValidator.VerifyChallenge(nonce, challengeParams.Challenges[0], challengeHash))
	})
}

func TestService_SetDifficulty(t *testing.T) {
	serviceHasher := hasher.NewService()
	appLog, err := logger.NewAppLogger("test")
	require.NoError(t, err)
	serviceValidator := validator.NewService(appLog, difficulty, maxAllowed, 3*time.Second, serviceHasher, nil)
	t.Run("change difficulty sha256", func(t *testing.T) {
		// given
		challengeParamsBefore := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParamsBefore.Difficulty, challengeParamsBefore.Challenges[0])
		require.True(t, serviceValidator.VerifyChallenge(nonce, challengeParamsBefore.Challenges[0], challengeHash))
		// when
		require.NoError(t, serviceValidator.SetDifficulty(challengeParamsBefore.Difficulty+1, entites.SHA256))

		// then
		challengeParamsAfter := serviceValidator.GetChallenges(1)
		challengeHashAfterWrong, nonceAfterWrong := getHashSHA256(serviceHasher, challengeParamsBefore.Difficulty, challengeParamsAfter.Challenges[0])
		require.False(t, serviceValidator.VerifyChallenge(nonceAfterWrong, challengeParamsAfter.Challenges[0], challengeHashAfterWrong))
		challengeHashAfter, nonceAfter := getHashSHA256(serviceHasher, challengeParamsAfter.Difficulty, challengeParamsAfter.Challenges[0])
		require.True(t, serviceValidator.VerifyChallenge(nonceAfter, challengeParamsAfter.Challenges[0], challengeHashAfter))
	})
	t.Run("change difficulty scrypt", func(t *testing.T) {
		// given
		challengeParamsBefore := serviceValidator.GetChallenges(1)
		challengeHash, nonce := getHashSHA256(serviceHasher, challengeParamsBefore.Difficulty, challengeParamsBefore.Challenges[0])
		require.True(t, serviceValidator.VerifyChallenge(nonce, challengeParamsBefore.Challenges[0], challengeHash))
		// when
		require.NoError(t, serviceValidator.SetDifficulty(challengeParamsBefore.Difficulty+1, entites.Scrypt))
		// then
		challengeParamsAfter := serviceValidator.GetChallenges(1)
		challengeHashAfter, nonceAfter := getHashScrypt(serviceHasher, challengeParamsAfter.Difficulty, challengeParamsAfter.Challenges[0])
		require.True(t, serviceValidator.VerifyChallenge(nonceAfter, challengeParamsAfter.Challenges[0], challengeHashAfter))
	})
}

func getHashSHA256(src *hasher.Service, difficulty uint32, base string) (hash string, nonce uint32) {
	compare := strings.Repeat("0", int(difficulty))
	for {
		h := src.EncodeSHA256(nonce, base)
		result := hex.EncodeToString(h)
		if result[:difficulty] == compare {
			return result, nonce
		}
		nonce++
	}
}

func getHashScrypt(src *hasher.Service, difficulty uint32, base string) (hash string, nonce uint32) {
	compare := strings.Repeat("0", int(difficulty))
	for {
		h := src.EncodeScrypt(nonce, base)
		result := hex.EncodeToString(h)
		if result[:difficulty] == compare {
			return result, nonce
		}
		nonce++
	}
}
