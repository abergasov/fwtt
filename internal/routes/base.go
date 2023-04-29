package routes

import (
	"fwtt/internal/logger"
	"fwtt/internal/service/quoter"
	"fwtt/internal/service/validator"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type Server struct {
	appAddr          string
	log              logger.AppLogger
	serviceQuoter    *quoter.Service
	serviceValidator *validator.Service
	httpEngine       *fiber.App
}

// InitAppRouter initializes the HTTP Server.
func InitAppRouter(log logger.AppLogger, serviceQuoter *quoter.Service, serviceValidator *validator.Service, address string) *Server {
	app := &Server{
		log:              log.With(zap.String("service", "http")),
		appAddr:          address,
		httpEngine:       fiber.New(fiber.Config{}),
		serviceQuoter:    serviceQuoter,
		serviceValidator: serviceValidator,
	}
	app.httpEngine.Use(recover.New())
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.httpEngine.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})
	s.httpEngine.Get("/challenges", func(ctx *fiber.Ctx) error {
		return ctx.JSON(s.serviceValidator.GetChallenges(ctx.QueryInt("num")))
	})
	s.httpEngine.Use(func(ctx *fiber.Ctx) error {
		challenge := ctx.Get("X-Challenge")
		if challenge == "" {
			return ctx.Status(http.StatusForbidden).SendString("challenge missing")
		}
		verifyHash := ctx.Get("X-Solution")
		if verifyHash == "" {
			return ctx.Status(http.StatusForbidden).SendString("solution missing")
		}
		nonceStr := ctx.Get("X-Nonce")
		if nonceStr == "" {
			return ctx.Status(http.StatusForbidden).SendString("nonce missing")
		}
		nonce, err := strconv.Atoi(nonceStr)
		if err != nil {
			return ctx.Status(http.StatusForbidden).SendString("nonce invalid")
		}
		if nonce <= 0 {
			return ctx.Status(http.StatusForbidden).SendString("nonce invalid")
		}
		if s.serviceValidator.VerifyChallenge(uint32(nonce), challenge, verifyHash) {
			return ctx.Next()
		}
		return ctx.Status(http.StatusForbidden).SendString("invalid solution")
	})
	s.httpEngine.Get("/quote", func(ctx *fiber.Ctx) error {
		return ctx.SendString(s.serviceQuoter.GetRandomQuote().Quote)
	})
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server", zap.String("port", s.appAddr))
	return s.httpEngine.Listen(s.appAddr)
}

func (s *Server) Stop() error {
	return s.httpEngine.Shutdown()
}
