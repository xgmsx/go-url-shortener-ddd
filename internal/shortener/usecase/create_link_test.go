package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
)

func TestCreateLink(t *testing.T) {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	tests := []struct {
		name        string
		database    []entity.Link
		databaseErr error
		cache       []entity.Link
		cacheErr    error
		broker      []entity.Link
		brokerErr   error
		input       dto.CreateLinkInput
		expected    dto.CreateLinkOutput
		expectedErr bool
	}{
		{
			name:     "link created",
			input:    dto.CreateLinkInput{URL: "http://example.com"},
			expected: dto.CreateLinkOutput{URL: "http://example.com"},
		},
		{
			name:        "link already exists",
			database:    []entity.Link{{URL: "http://example.com", Alias: "test_link"}},
			input:       dto.CreateLinkInput{URL: "http://example.com"},
			expected:    dto.CreateLinkOutput{URL: "http://example.com", Alias: "test_link"},
			expectedErr: true,
		},
		{
			name:        "link created, but error occurred in cache",
			cacheErr:    fmt.Errorf("cache timeout"),
			input:       dto.CreateLinkInput{URL: "http://example.com"},
			expected:    dto.CreateLinkOutput{URL: "http://example.com"},
			expectedErr: true,
		},
		{
			name:        "link created, but error occurred in broker",
			brokerErr:   fmt.Errorf("broker timeout"),
			input:       dto.CreateLinkInput{URL: "http://example.com"},
			expected:    dto.CreateLinkOutput{URL: "http://example.com"},
			expectedErr: true,
		},
		{
			name:        "error in database",
			databaseErr: fmt.Errorf("database timeout"),
			input:       dto.CreateLinkInput{URL: "http://example.com"},
			expected:    dto.CreateLinkOutput{},
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
			output, err := usecase.New(db, cache, broker).CreateLink(ctx, tt.input)

			// Assert
			switch {
			case tt.expectedErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.expected.URL, output.URL)
				assert.NotEmpty(t, output.Alias)
			}
		})
	}
}
