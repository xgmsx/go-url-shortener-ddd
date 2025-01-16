package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"

	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/adapter/redis"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
)

const ttl = time.Hour

func marshal(link entity.Link) []byte {
	d, _ := json.Marshal(link)
	return d
}

func newFakeRedisClient() (*redis.Redis, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	return redis.New(client), mock
}

func TestPutLink(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		want    entity.Link
		wantErr error
	}{
		{
			name:  "happy path",
			alias: "test_alias",
			want:  entity.Link{URL: "http://example.com", Alias: "test_alias"},
		},
		{
			name:    "error path",
			alias:   "test_alias",
			want:    entity.Link{URL: "http://example.com", Alias: "test_alias"},
			wantErr: fmt.Errorf("test redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			r, mock := newFakeRedisClient()
			if tt.wantErr != nil {
				mock.ExpectSet(tt.alias, marshal(tt.want), ttl).SetErr(tt.wantErr)
			} else {
				mock.ExpectSet(tt.alias, marshal(tt.want), ttl).SetVal("OK")
			}

			// Act
			err := r.PutLink(context.Background(), tt.want)

			// Assert
			require.NoError(t, mock.ExpectationsWereMet())
			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetLink(t *testing.T) {
	tests := []struct {
		name    string
		alias   string
		want    entity.Link
		wantErr error
	}{
		{
			name:  "happy path",
			alias: "test_alias",
			want:  entity.Link{URL: "http://example.com", Alias: "test_alias"},
		},
		{
			name:    "not found path",
			alias:   "nonexistent_alias",
			want:    entity.Link{},
			wantErr: entity.ErrNotFound,
		},
		{
			name:    "error path",
			alias:   "test_alias",
			want:    entity.Link{URL: "http://example.com", Alias: "test_alias"},
			wantErr: fmt.Errorf("test redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			r, mock := newFakeRedisClient()
			if tt.wantErr != nil {
				mock.ExpectGet(tt.alias).SetErr(tt.wantErr)
			} else {
				mock.ExpectGet(tt.alias).SetVal(string(marshal(tt.want)))
			}

			// Act
			got, err := r.GetLink(context.Background(), tt.alias)

			// Assert
			require.NoError(t, mock.ExpectationsWereMet())
			if tt.wantErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
