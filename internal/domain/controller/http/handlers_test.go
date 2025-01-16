package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	ucCreate "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	mocksCreate "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create/mocks"
	ucFetch "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
	mocksFetch "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch/mocks"
)

func TestCreateLink(t *testing.T) {
	initMock := func() (*gomock.Controller, *mocksCreate.Mockdatabase, *mocksCreate.Mockcache, *mocksCreate.Mockpublisher) {
		ctrl := gomock.NewController(t)
		publisher := mocksCreate.NewMockpublisher(ctrl)
		database := mocksCreate.NewMockdatabase(ctrl)
		cache := mocksCreate.NewMockcache(ctrl)
		return ctrl, database, cache, publisher
	}

	testCases := []struct {
		name       string
		input      string
		wantStatus int
		wantOutput string
		setupMock  func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher)
	}{
		{
			name:       "Happy path",
			input:      `{"url": "https://example.com"}`,
			wantStatus: http.StatusCreated,
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				cache.EXPECT().PutLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				publisher.EXPECT().SendLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},

		{
			name:       "Invalid JSON",
			input:      `test text`,
			wantStatus: http.StatusBadRequest,
			wantOutput: "invalid json",
		},

		{
			name:       "JSON validation error",
			input:      `{"test": "json"}`,
			wantStatus: http.StatusBadRequest,
			wantOutput: "validation error",
		},
		{
			name:       "Link already exists",
			input:      `{"url": "https://example.com"}`,
			wantStatus: http.StatusFound,
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(entity.ErrAlreadyExist).Times(1)
			},
		},
		{
			name:       "Internal error",
			input:      `{"url": "https://example.com"}`,
			wantStatus: http.StatusInternalServerError,
			wantOutput: "internal error",
			setupMock: func(database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(errors.New("test error")).Times(1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, database, cache, publisher := initMock()
			defer ctrl.Finish()

			if tc.setupMock != nil {
				tc.setupMock(database, cache, publisher)
			}

			// arrange
			uc := ucCreate.New(database, cache, publisher)

			srv := fiber.New()
			srv.Add(http.MethodPost, "/create", NewHandlerCreateLink(uc).Handler)

			// act
			resp, output := sendHTTPRequest(t, srv, http.MethodPost, "/create", tc.input)

			// assert
			if tc.wantStatus != 0 {
				assert.Equal(t, tc.wantStatus, resp.StatusCode)
			}
			if tc.wantOutput != "" {
				assert.Equal(t, tc.wantOutput, output)
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
		alias      string
		wantStatus int
		wantOutput string
		setupMock  func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache)
	}{
		{
			name:       "Happy path",
			alias:      "alias1",
			wantStatus: http.StatusOK,
			wantOutput: `{"url":"https://example.com","alias":"alias1","expired_at":"0001-01-01T00:00:00Z"}`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				link := entity.Link{URL: "https://example.com", Alias: "alias1", ExpiredAt: time.Time{}}
				cache.EXPECT().GetLink(gomock.Any(), "alias1").Return(nil, entity.ErrNotFound).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "alias1", "").Return(&link, nil).Times(1)
				cache.EXPECT().PutLink(gomock.Any(), link).Return(nil).Times(1)
			},
		},
		{
			name:       "Happy path with cached link",
			alias:      `alias2`,
			wantStatus: http.StatusOK,
			wantOutput: `{"url":"https://example.com","alias":"alias2","expired_at":"0001-01-01T00:00:00Z"}`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cachedLink := entity.Link{URL: "https://example.com", Alias: "alias2", ExpiredAt: time.Time{}}
				cache.EXPECT().GetLink(gomock.Any(), "alias2").Return(&cachedLink, nil).Times(1)
			},
		},
		{
			name:       "Link not found",
			alias:      "unknown",
			wantStatus: http.StatusNotFound,
			wantOutput: `not found`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cache.EXPECT().GetLink(gomock.Any(), "unknown").Return(nil, entity.ErrNotFound).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "unknown", "").Return(nil, entity.ErrNotFound).Times(1)
			},
		},
		{
			name:       "Internal error",
			alias:      "alias3",
			wantStatus: http.StatusInternalServerError,
			wantOutput: `internal error`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cache.EXPECT().GetLink(gomock.Any(), "alias3").Return(nil, errors.New("test cache error")).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "alias3", "").Return(nil, errors.New("test db error")).Times(1)
			},
		},
		{
			name:       "Validation error",
			alias:      "a",
			wantStatus: http.StatusBadRequest,
			wantOutput: `validation error`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl, database, cache := initMock()
			defer ctrl.Finish()

			if tc.setupMock != nil {
				tc.setupMock(database, cache)
			}

			// arrange
			uc := ucFetch.New(database, cache)

			srv := fiber.New()
			srv.Add(http.MethodGet, "/fetch/:alias", NewHandlerFetchLink(uc).Handler)

			// act
			resp, output := sendHTTPRequest(t, srv, http.MethodGet, "/fetch/"+tc.alias, "")

			// assert
			if tc.wantStatus != 0 {
				assert.Equal(t, tc.wantStatus, resp.StatusCode)
			}
			if tc.wantOutput != "" {
				assert.Equal(t, tc.wantOutput, output)
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	initMock := func() (*gomock.Controller, *mocksFetch.Mockdatabase, *mocksFetch.Mockcache) {
		ctrl := gomock.NewController(t)
		database := mocksFetch.NewMockdatabase(ctrl)
		cache := mocksFetch.NewMockcache(ctrl)
		return ctrl, database, cache
	}

	testCases := []struct {
		name       string
		alias      string
		wantStatus int
		wantOutput string
		setupMock  func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache)
	}{
		{
			name:       "Happy path",
			alias:      "alias1",
			wantStatus: http.StatusFound,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				link := entity.Link{URL: "https://example.com", Alias: "alias1", ExpiredAt: time.Time{}}
				cache.EXPECT().GetLink(gomock.Any(), "alias1").Return(&link, nil).Times(1)
			},
		},
		{
			name:       "Link not found",
			alias:      "unknown",
			wantStatus: http.StatusNotFound,
			wantOutput: `not found`,
			setupMock: func(database *mocksFetch.Mockdatabase, cache *mocksFetch.Mockcache) {
				cache.EXPECT().GetLink(gomock.Any(), "unknown").Return(nil, entity.ErrNotFound).Times(1)
				database.EXPECT().FindLink(gomock.Any(), "unknown", "").Return(nil, entity.ErrNotFound).Times(1)
			},
		},
		{
			name:       "Validation error",
			alias:      "a",
			wantStatus: http.StatusBadRequest,
			wantOutput: `validation error`,
		},
		{
			name:       "Internal error",
			alias:      "alias3",
			wantStatus: http.StatusInternalServerError,
			wantOutput: `internal error`,
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

			// arrange
			uc := ucFetch.New(database, cache)

			srv := fiber.New()
			srv.Add(http.MethodGet, "/fetch/:alias/redirect", NewHandlerRedirect(uc).Handler)

			// act
			resp, output := sendHTTPRequest(t, srv, http.MethodGet, "/fetch/"+tc.alias+"/redirect", "")

			// assert
			if tc.wantStatus != 0 {
				assert.Equal(t, tc.wantStatus, resp.StatusCode)
			}
			if tc.wantOutput != "" {
				assert.Equal(t, tc.wantOutput, output)
			}
		})
	}
}

func sendHTTPRequest(test *testing.T, app *fiber.App, method, url, body string) (resp *http.Response, respBody string) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(test, err)
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(test, err)

	return resp, string(respBodyBytes)
}
