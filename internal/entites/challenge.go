package entites

import (
	"sync"
	"sync/atomic"
)

type Challenge struct {
	Challenge  string `db:"challenge"`
	ValidTill  int64  `db:"valid_till"`
	Difficulty uint32 `db:"difficulty"`
	MaxAllowed uint32 `db:"max_allowed"`
	Used       uint32 `db:"used"`
	HashAlgo   string `db:"hash_algo"`
	Hash       string `db:"hash"`
	mu         sync.RWMutex
}

func (c *Challenge) SetHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Hash = hash
}

func (c *Challenge) IncrementUsage() {
	atomic.AddUint32(&c.Used, 1)
}

func (c *Challenge) GetHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Hash
}
