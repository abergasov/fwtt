package validator

import (
	"errors"
	"fmt"
	"fwtt/internal/entites"
	"fwtt/internal/logger"
	"fwtt/internal/repository/validator"
	"fwtt/internal/service/hasher"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrUnknownHash = errors.New("unknown hash algorithm")
)

// Service is a validator service. It generates challenges and validates solutions
type Service struct {
	difficulty     uint32
	maxDuration    time.Duration
	maxAllowed     uint32
	challenges     map[string]*entites.Challenge
	challengesAlgo string
	challengesMU   sync.RWMutex

	log           logger.AppLogger
	serviceHasher *hasher.Service
	repoValidator *validator.Repository
}

func NewService(log logger.AppLogger, difficulty, maxAllowed uint32, maxDuration time.Duration, serviceHasher *hasher.Service, repoValidator *validator.Repository) *Service {
	srv := &Service{
		difficulty:     difficulty,
		maxAllowed:     maxAllowed,
		maxDuration:    maxDuration,
		challengesAlgo: entites.SHA256,
		challenges:     make(map[string]*entites.Challenge),

		log:           log.With(zap.String("service", "validator")),
		serviceHasher: serviceHasher,
		repoValidator: repoValidator,
	}
	if err := srv.loadState(); err != nil {
		srv.log.Fatal("failed to load state", err)
	}
	go srv.emulateStateObserver()
	return srv
}

// emulateStateObserver is a helper function that changes challengesAlgo every 5 seconds
func (s *Service) emulateStateObserver() {
	for range time.NewTicker(5 * time.Second).C {
		s.challengesMU.Lock()
		if s.challengesAlgo == entites.SHA256 {
			s.challengesAlgo = entites.Scrypt
		} else {
			s.challengesAlgo = entites.SHA256
		}
		s.log.Info("changed challenges algo", zap.String("algo", s.challengesAlgo))
		s.challengesMU.Unlock()
	}
}

func (s *Service) loadState() error {
	if s.repoValidator == nil {
		return nil
	}
	challengesList, err := s.repoValidator.LoadChallenges()
	if err != nil {
		return fmt.Errorf("failed to load challenges: %w", err)
	}
	removeChallenges := make([]string, 0, len(challengesList))
	now := time.Now().Unix()
	for _, ch := range challengesList {
		if ch.Used >= ch.MaxAllowed || now > ch.ValidTill {
			removeChallenges = append(removeChallenges, ch.Challenge)
			continue
		}
		s.challenges[ch.Challenge] = ch
	}
	if err = s.repoValidator.DropChallenges(removeChallenges); err != nil {
		return fmt.Errorf("failed to drop challenges: %w", err)
	}
	return nil
}

func (s *Service) Stop() {
	if s.repoValidator == nil {
		return
	}
	s.challengesMU.Lock()
	defer s.challengesMU.Unlock()
	challenges := make([]*entites.Challenge, 0, len(s.challenges))
	for _, ch := range s.challenges {
		if ch.Used >= ch.MaxAllowed || time.Now().Unix() > ch.ValidTill {
			continue
		}
		challenges = append(challenges, ch)
	}
	if err := s.repoValidator.SaveChallenges(challenges); err != nil {
		s.log.Error("failed to save challenges", err)
	}
}

func (s *Service) SetDifficulty(difficulty uint32, algo string) error {
	if !entites.ValidateHashAlgo(algo) {
		return ErrUnknownHash
	}
	atomic.StoreUint32(&s.difficulty, difficulty)
	s.challengesMU.Lock()
	s.challengesAlgo = algo
	s.challengesMU.Unlock()
	return nil
}

func (s *Service) GetDifficulty() uint32 {
	return atomic.LoadUint32(&s.difficulty)
}

// GetChallenges returns 10 challenges
func (s *Service) GetChallenges(num int) *entites.Challenges {
	if num <= 0 {
		return nil
	}
	if num > 10 {
		num = 10
	}
	challenges := make([]string, num)
	currentDiff := atomic.LoadUint32(&s.difficulty)
	s.challengesMU.Lock()
	defer s.challengesMU.Unlock()
	for i := range challenges {
		ch := uuid.NewString()
		challenges[i] = ch
		s.challenges[ch] = &entites.Challenge{
			ValidTill:  time.Now().Add(s.maxDuration).Unix(),
			Challenge:  ch,
			Difficulty: currentDiff,
			MaxAllowed: s.maxAllowed,
			HashAlgo:   s.challengesAlgo,
		}
	}
	result := &entites.Challenges{
		Challenges: challenges,
		Difficulty: currentDiff,
		Algorithm:  s.challengesAlgo,
	}
	if s.challengesAlgo == entites.Scrypt {
		result.AlgoParams = s.serviceHasher.ScryptConfig()
	}
	return result
}

// VerifyChallenge verifies the challenge and returns true if it is valid
// logic:
// 1. check if challenge exists
// 2. check if challenge is not expired
// 3. check if challenge is not used
// 4. check if challenge is valid
func (s *Service) VerifyChallenge(nonce uint32, challenge, verifyHash string) bool {
	s.log.Info("verify challenge", zap.String("challenge", challenge), zap.Uint32("nonce", nonce), zap.String("verifyHash", verifyHash))
	s.challengesMU.RLock()
	ch, ok := s.challenges[challenge]
	s.challengesMU.RUnlock()
	if !ok { // unknown challenge
		return false
	}
	if ch.Used >= ch.MaxAllowed { // challenge already used
		return false
	}
	if time.Now().Unix() > ch.ValidTill { // challenge expired
		return false
	}
	knownHash := ch.GetHash()
	if knownHash != "" && knownHash != verifyHash { // challenge is not valid
		return false
	}
	// check if challenge is valid
	var hashValid bool
	if ch.HashAlgo == entites.SHA256 {
		hashValid = s.serviceHasher.CheckSHA256(ch.Difficulty, nonce, ch.Challenge, verifyHash)
	} else if ch.HashAlgo == entites.Scrypt {
		hashValid = s.serviceHasher.CheckScrypt(ch.Difficulty, nonce, ch.Challenge, verifyHash)
	} else {
		s.log.Error("error check solution", ErrUnknownHash, zap.String("algo", ch.HashAlgo))
		return false
	}
	if hashValid {
		if ch.Used+1 > ch.MaxAllowed {
			// remove challenge from map
			s.challengesMU.Lock()
			delete(s.challenges, challenge)
			s.challengesMU.Unlock()
		}
		ch.SetHash(verifyHash)
		atomic.AddUint32(&ch.Used, 1)
		return true
	}
	return false
}
