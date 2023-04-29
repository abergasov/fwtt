package hasher_test

import (
	"encoding/hex"
	"fwtt/internal/service/hasher"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const (
	difficulty = 3
	start      = "000"
)

func TestService_CheckSHA256(t *testing.T) {
	srv := hasher.NewService()
	hash, base, nonce := generateSHA(srv)
	t.Log("hash:", hash)
	t.Log("base:", base)
	t.Log("nonce:", nonce)
	require.True(t, srv.CheckSHA256(difficulty, nonce, base, hash))
}

func TestService_CheckScrypt(t *testing.T) {
	srv := hasher.NewService()
	hash, base, nonce := generateScrypt(srv)
	t.Log("hash:", hash)
	t.Log("base:", base)
	t.Log("nonce:", nonce)
	require.True(t, srv.CheckScrypt(difficulty, nonce, base, hash))
}

func BenchmarkService_CheckSHA256(b *testing.B) {
	srv := hasher.NewService()
	hash, base, nonce := generateSHA(srv)
	for i := 0; i < b.N; i++ {
		require.True(b, srv.CheckSHA256(difficulty, nonce, base, hash))
	}
}

func BenchmarkService_CheckScrypt(b *testing.B) {
	srv := hasher.NewService()
	hash, base, nonce := generateScrypt(srv)
	for i := 0; i < b.N; i++ {
		require.True(b, srv.CheckScrypt(difficulty, nonce, base, hash))
	}
}

func generateSHA(srv *hasher.Service) (hash, challenge string, nonce uint32) {
	base := uuid.NewString()
	for {
		res := hex.EncodeToString(srv.EncodeSHA256(nonce, base))
		if res[:difficulty] == start {
			return res, base, nonce
		}
		nonce++
	}
}

func generateScrypt(srv *hasher.Service) (hash, challenge string, nonce uint32) {
	base := uuid.NewString()
	for {
		res := hex.EncodeToString(srv.EncodeScrypt(nonce, base))
		if res[:difficulty] == start {
			return res, base, nonce
		}
		nonce++
	}
}
