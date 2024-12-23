package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const readTimeout = 5 * time.Second
const writeTimeout = 5 * time.Second
const closeTimeout = 5 * time.Second

type Option func(*Options)

type Registrable interface {
	Register(prefix string, app *fiber.App)
}

type Options struct {
	Routers      []map[string]Registrable
	Middlewares  []fiber.Handler
	Name         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
}

func BuildOptions(c Config, opts ...Option) *Options {
	var o = &Options{
		Name:         c.AppName,
		Port:         c.Port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		CloseTimeout: closeTimeout,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithName(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func WithPort(port string) Option {
	return func(options *Options) {
		options.Port = port
	}
}

func WithCloseTimeout(d time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = d
	}
}

func WithReadTimeout(d time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(options *Options) {
		options.CloseTimeout = d
	}
}

func WithRouter(prefix string, router Registrable) Option {
	return func(options *Options) {
		options.Routers = append(options.Routers, map[string]Registrable{prefix: router})
	}
}

func WithMiddleware(m fiber.Handler) Option {
	return func(options *Options) {
		options.Middlewares = append(options.Middlewares, m)
	}
}

func WithMiddlewares(m []fiber.Handler) Option {
	return func(options *Options) {
		options.Middlewares = m
	}
}
