package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel"
)

func TestGetLink(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	err := otel.Init(ctx, otel.Config{}, "test app", "0.0.0")
	require.NoError(t, err)

	tests := []struct {
		name        string
		database    []entity.Link
		databaseErr error
		cache       []entity.Link
		cacheErr    error
		broker      []entity.Link
		brokerErr   error
		input       dto.GetLinkInput
		expected    dto.GetLinkOutput
		expectedErr bool
	}{
		{
			name:     "link found in cache",
			cache:    []entity.Link{{Alias: "cache_link", URL: "http://example1.com"}},
			database: []entity.Link{{Alias: "db_link", URL: "http://example2.com"}},
			input:    dto.GetLinkInput{Alias: "cache_link"},
			expected: dto.GetLinkOutput{Alias: "cache_link", URL: "http://example1.com"},
		},
		{
			name:     "link found in database",
			cache:    []entity.Link{{Alias: "cache_link", URL: "http://example3.com"}},
			database: []entity.Link{{Alias: "db_link", URL: "http://example4.com"}},
			input:    dto.GetLinkInput{Alias: "db_link"},
			expected: dto.GetLinkOutput{Alias: "db_link", URL: "http://example4.com"},
		},
		{
			name:        "link not found in cache or database",
			cache:       []entity.Link{{Alias: "test_link1", URL: "http://example5.com"}},
			database:    []entity.Link{{Alias: "test_lint1", URL: "http://example6.com"}},
			input:       dto.GetLinkInput{Alias: "not_existing_link"},
			expectedErr: true,
		},
		{
			name:     "link found in database, but error occurred in cache",
			database: []entity.Link{{Alias: "test_lint2", URL: "http://example6.com"}},
			cache:    []entity.Link{{Alias: "test_lint2", URL: "http://example7.com"}},
			cacheErr: fmt.Errorf("timeout"),
			input:    dto.GetLinkInput{Alias: "test_lint2"},
			expected: dto.GetLinkOutput{Alias: "test_lint2", URL: "http://example6.com"},
		},
		{
			name:        "database error while getting link",
			databaseErr: fmt.Errorf("database timeout"),
			input:       dto.GetLinkInput{Alias: "db_link"},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			db := newFakeDatabase(tt.database)
			db.Err = tt.databaseErr
			cache := newFakeCache(tt.cache)
			cache.Err = tt.cacheErr
			broker := newFakeBroker(tt.broker)
			broker.Err = tt.brokerErr

			// Act
			output, err := usecase.New(db, cache, broker).GetLink(ctx, tt.input)

			// Assert
			switch {
			case tt.expectedErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.expected, output)
				assert.WithinDuration(t, tt.expected.ExpiredAt, output.ExpiredAt, time.Second)
			}
		})
	}
}
