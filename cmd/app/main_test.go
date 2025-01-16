package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
)

func TestRun(t *testing.T) {
	envs := map[string]string{
		"APP_NAME":           "test",
		"APP_VERSION":        "0.0.0",
		"POSTGRES_USER":      "testuser",
		"POSTGRES_PASSWORD":  "testpass",
		"POSTGRES_HOST":      "testhost",
		"POSTGRES_DB":        "testdb",
		"REDIS_ADDR":         "testredis",
		"KAFKA_BROKERS":      "testkafka",
		"KAFKA_OUTPUT_TOPIC": "topic1",
		"KAFKA_INPUT_TOPIC":  "topic2",
		"KAFKA_GROUP":        "test",
	}
	for key, value := range envs {
		t.Setenv(key, value)
	}

	tests := []struct {
		name string
		fn   func(context.Context, *config.Config) error
	}{
		{
			name: "happy path",
			fn:   func(_ context.Context, _ *config.Config) error { return nil },
		},
		{
			name: "error path",
			fn:   func(_ context.Context, _ *config.Config) error { return fmt.Errorf("test error") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			run(tt.fn)
		})
	}
}
