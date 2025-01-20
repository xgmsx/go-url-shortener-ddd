package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"

	"github.com/xgmsx/go-url-shortener-ddd/pkg/http/middlewares/metrics"

	"github.com/xgmsx/go-url-shortener-ddd/docs"
)

type Config struct {
	AppName      string `env:"APP_NAME, required"`
	AppVersion   string `env:"APP_VERSION, required"`
	Port         string `env:"HTTP_PORT, default=8000"`
	UseRecover   bool   `env:"HTTP_USE_RECOVER, default=true"`
	UseRequestID bool   `env:"HTTP_USE_REQUEST_ID, default=true"`
	UsePprof     bool   `env:"HTTP_USE_PPROF, default=false"`
}

type registrable interface {
	Register(s *fiber.App)
}

// Server HTTP.
//
// @version      0.0.0
// @title        Title
// @BasePath     /api
type Server struct {
	app     *fiber.App
	options *Options
	config  Config
}

func New(config Config, options *Options, controllers ...registrable) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler:          DefaultErrorHandler,
		AppName:               options.Name,
		ReadTimeout:           options.ReadTimeout,
		WriteTimeout:          options.WriteTimeout,
		DisableStartupMessage: true,
	})

	// middlewares
	for _, m := range options.Middlewares {
		app.Use(m)
	}

	// controllers
	for _, c := range controllers {
		c.Register(app)
	}

	app.Get("/live", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/ready", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/metrics", metrics.HandlerDefault)

	docs.SwaggerInfo.Title = options.Name
	docs.SwaggerInfo.Version = config.AppVersion

	return &Server{
		app:     app,
		config:  config,
		options: options,
	}
}

func (s *Server) Serve(port string) error {
	log.Info().Msg("HTTP server started on port: " + port)
	return s.app.Listen("0.0.0.0:" + port)
}

func (s *Server) Close() {
	err := s.app.ShutdownWithTimeout(s.options.CloseTimeout)
	if err != nil {
		log.Error().Err(err).Msg("server - Close - s.app.ShutdownWithTimeout")
	}

	log.Info().Msg("HTTP server closed")
}
