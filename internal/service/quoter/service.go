package quoter

import (
	"errors"
	"fmt"
	"fwtt/internal/entites"
	"fwtt/internal/logger"
	"fwtt/internal/repository/quotes"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	rand       int64
	quotesLen  int64
	log        logger.AppLogger
	repo       *quotes.Repository
	quotesList []*entites.Quote
}

func NewService(log logger.AppLogger, repo *quotes.Repository) *Service {
	source := rand.NewSource(time.Now().UnixNano())
	srv := &Service{
		rand: rand.New(source).Int63(), // nolint:gosec
		repo: repo,
		log:  log.With(zap.String("service", "qouter")),
	}
	if err := srv.loadQuotes(); err != nil {
		log.Fatal("can't load quotes", err)
	}
	return srv
}

func (s *Service) loadQuotes() error {
	quotesList, err := s.repo.LoadQuotes()
	if err != nil {
		return fmt.Errorf("failed to load quotes: %w", err)
	}
	s.quotesList = quotesList
	s.quotesLen = int64(len(quotesList))
	if s.quotesLen == 0 {
		return errors.New("quotes list is empty")
	}
	return nil
}

func (s *Service) GetRandomQuote() *entites.Quote {
	if s.quotesLen == 0 {
		return &entites.Quote{}
	}
	index := (time.Now().UnixNano() + s.rand) % s.quotesLen
	if index < 0 || index >= s.quotesLen {
		index = 0
	}
	return s.quotesList[index]
}
