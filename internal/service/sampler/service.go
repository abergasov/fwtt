package sampler

import (
	"fwtt/internal/logger"
	"fwtt/internal/repository/sampler"

	"go.uber.org/zap"
)

type Service struct {
	log  logger.AppLogger
	repo *sampler.Repo
}

func InitService(log logger.AppLogger, repo *sampler.Repo) *Service {
	return &Service{
		repo: repo,
		log:  log.With(zap.String("service", "sampler")),
	}
}
