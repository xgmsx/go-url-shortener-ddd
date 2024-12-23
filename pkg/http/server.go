package http

import (
	"url-shortener/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port       string `env:"HTTP_PORT, default=8000"`
	AppName    string `env:"APP_NAME, required"`
	AppVersion string `env:"APP_VERSION, required"`
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
	notify  chan error
}

func New(ch chan error, config Config, options *Options) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler:          DefaultErrorHandler,
		AppName:               options.Name,
		ReadTimeout:           options.ReadTimeout,
		WriteTimeout:          options.WriteTimeout,
		DisableStartupMessage: true,
	})

	s := &Server{
		app:     app,
		config:  config,
		options: options,
		notify:  ch,
	}

	docs.SwaggerInfo.Title = options.Name
	docs.SwaggerInfo.Version = config.AppVersion

	// middlewares
	for _, middleware := range options.Middlewares {
		app.Use(middleware)
	}

	// routers
	for _, r := range options.Routers {
		for prefix, router := range r {
			router.Register(prefix, app)
		}
	}

	go func() {
		s.Notify(s.app.Listen("0.0.0.0:" + config.Port))
	}()

	log.Info().Msg("HTTP server started on port: " + config.Port)

	return s
}

func (s *Server) Close() {
	err := s.app.ShutdownWithTimeout(s.options.CloseTimeout)
	if err != nil {
		log.Error().Err(err).Msg("server - Close - s.app.ShutdownWithTimeout")
	}

	log.Info().Msg("HTTP server closed")
}

func (s *Server) Notify(err error) {
	if err != nil {
		s.notify <- err
	}
}
