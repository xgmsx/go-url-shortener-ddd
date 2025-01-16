package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/controller/grpc"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	mocksCreate "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create/mocks"
	ucFetch "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
	mocksFetch "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch/mocks"
	pb "github.com/xgmsx/go-url-shortener-ddd/proto/gen/shortener.v1"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateLink(t *testing.T) {
	initMock := func(t *testing.T) (*gomock.Controller, *mocksCreate.Mockdatabase, *mocksCreate.Mockcache, *mocksCreate.Mockpublisher) {
		ctrl := gomock.NewController(t)
		database := mocksCreate.NewMockdatabase(ctrl)
		cache := mocksCreate.NewMockcache(ctrl)
		publisher := mocksCreate.NewMockpublisher(ctrl)
		return ctrl, database, cache, publisher
	}

	testCases := []struct {
		name       string
		input      *pb.CreateLinkRequest
		wantStatus codes.Code
		wantOutput *pb.CreateLinkResponse
		wantError  string
		setupMock  func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher)
	}{
		{
			name:       "Happy path",
			input:      &pb.CreateLinkRequest{Url: "https://example.com"},
			wantStatus: codes.OK,
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				cache.EXPECT().PutLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				publisher.EXPECT().SendLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:       "Link already exists",
			input:      &pb.CreateLinkRequest{Url: "https://example.com"},
			wantStatus: codes.OK,
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(entity.ErrAlreadyExist).Times(1)
			},
		},
		{
			name:       "Validation error",
			input:      &pb.CreateLinkRequest{Url: ""},
			wantStatus: codes.OK,
			wantError:  `validation error`,
		},
		{
			name:       "Internal error",
			input:      &pb.CreateLinkRequest{Url: "https://example.com"},
			wantStatus: codes.OK,
			wantError:  "internal error",
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, database, cache, publisher := initMock(t)
			defer ctrl.Finish()

			if tc.setupMock != nil {
				tc.setupMock(database, cache, publisher)
			}

			// arrange
			uc := create.New(database, cache, publisher)
			handler := grpc.NewHandlerCreateLink(uc)

			// act
			resp, err := handler.CreateLink(context.Background(), tc.input)

			// assert
			if tc.wantStatus != codes.OK {
				assert.Nil(t, resp)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.wantStatus, st.Code())
			}

			if tc.wantOutput != nil {
				require.NoError(t, err)
				assert.Equal(t, tc.wantOutput.Url, resp.Url)
				assert.Equal(t, tc.wantOutput.Alias, resp.Alias)
				assert.Equal(t, tc.wantOutput.ExpiredAt.AsTime().Truncate(time.Second), resp.ExpiredAt.AsTime().Truncate(time.Second))
			}

			if tc.wantError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantError)
			}
		})
	}
}

func TestFetchLink(t *testing.T) {
	initMock := func() (*gomock.Controller, *mocksFetch.Mockdatabase, *mocksFetch.Mockcache) {
		ctrl := gomock.NewController(t)
		database := mocksFetch.NewMockdatabase(ctrl)
		cache := mocksFetch.NewMockcache(ctrl)
		return ctrl, database, cache
	}

	testCases := []struct {
		name       string
		input      *pb.FetchLinkRequest
		wantStatus codes.Code
		wantOutput *pb.FetchLinkResponse
		wantError  string
		setupMock  func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache)
	}{
		{
			name:       "Happy path",
			input:      &pb.FetchLinkRequest{Alias: "alias1"},
			wantStatus: codes.OK,
			wantOutput: &pb.FetchLinkResponse{
				Url: "https://example.com", Alias: "alias1", ExpiredAt: timestamppb.New(time.Time{}),
			},
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				link := entity.Link{URL: "https://example.com", Alias: "alias1", ExpiredAt: time.Time{}}
				cache.EXPECT().GetLink(gomock.Any(), "alias1").Return(nil, entity.ErrNotFound).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "alias1", "").Return(&link, nil).Times(1)
				cache.EXPECT().PutLink(gomock.Any(), link).Return(nil).Times(1)
			},
		},
		{
			name:       "Happy path with cached link",
			input:      &pb.FetchLinkRequest{Alias: "alias2"},
			wantStatus: codes.OK,
			wantOutput: &pb.FetchLinkResponse{
				Url: "https://example.com", Alias: "alias2", ExpiredAt: timestamppb.New(time.Time{}),
			},
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cachedLink := entity.Link{URL: "https://example.com", Alias: "alias2", ExpiredAt: time.Time{}}
				cache.EXPECT().GetLink(gomock.Any(), "alias2").Return(&cachedLink, nil).Times(1)
			},
		},
		{
			name:       "Link not found",
			input:      &pb.FetchLinkRequest{Alias: "unknown"},
			wantStatus: codes.OK,
			wantError:  `not found`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cache.EXPECT().GetLink(gomock.Any(), "unknown").Return(nil, entity.ErrNotFound).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "unknown", "").Return(nil, entity.ErrNotFound).Times(1)
			},
		},
		{
			name:       "Validation error",
			input:      &pb.FetchLinkRequest{Alias: "a"},
			wantStatus: codes.OK,
			wantError:  `validation error`,
		},
		{
			name:       "Internal error",
			input:      &pb.FetchLinkRequest{Alias: "alias3"},
			wantStatus: codes.OK,
			wantError:  `internal error`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cache.EXPECT().GetLink(gomock.Any(), "alias3").Return(nil, errors.New("test cache error")).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "alias3", "").Return(nil, errors.New("test db error")).Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, database, cache := initMock()
			defer ctrl.Finish()

			if tc.setupMock != nil {
				tc.setupMock(database, cache)
			}

			/// arrange
			uc := ucFetch.New(database, cache)
			handler := grpc.NewHandlerFetchLink(uc)

			// act
			resp, err := handler.FetchLink(context.Background(), tc.input)

			// assert
			if tc.wantStatus != codes.OK {
				assert.Nil(t, resp)
				st, _ := status.FromError(err)
				assert.Equal(t, tc.wantStatus, st.Code())
			}

			if tc.wantOutput != nil {
				require.NoError(t, err)
				assert.Equal(t, tc.wantOutput.Url, resp.Url)
				assert.Equal(t, tc.wantOutput.Alias, resp.Alias)
				assert.Equal(t, tc.wantOutput.ExpiredAt.AsTime().Truncate(time.Second), resp.ExpiredAt.AsTime().Truncate(time.Second))
			}

			if tc.wantError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantError)
			}
		})
	}
}
