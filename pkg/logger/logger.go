package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level         string `env:"LOGGER_LEVEL, default=error"`
	PrettyConsole bool   `env:"LOGGER_PRETTY_CONSOLE, default=false"`
}

func Init(c Config, name, version string) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	level, err := zerolog.ParseLevel(c.Level)
	if err != nil {
		zerolog.SetGlobalLevel(level)
	}

	if c.PrettyConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
	} else {
		log.Logger = log.With().
			Caller().
			Str("app_name", name).
			Str("app_version", version).
			Logger()
	}

	log.Info().Msg("Logger initialized")
}
