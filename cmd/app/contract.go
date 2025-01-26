package main

import (
	"context"

	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
)

//go:generate mockgen -source=contract.go -destination=mocks/contract.go

type appRunner interface {
	Run(ctx context.Context, c *config.Config) error
}

type configLoader interface {
	Load(ctx context.Context) (*config.Config, error)
}
