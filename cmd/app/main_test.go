// cmd/app/main_test.go
package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mocks "github.com/xgmsx/go-url-shortener-ddd/cmd/app/mocks"
	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		mockConfig func(loader *mocks.MockconfigLoader)
		mockApp    func(runner *mocks.MockappRunner)
		wantErr    string
	}{
		{
			name: "Happy path",
			mockConfig: func(cl *mocks.MockconfigLoader) {
				cl.EXPECT().Load(gomock.Any()).Return(config.New(), nil).Times(1)
			},
			mockApp: func(ar *mocks.MockappRunner) {
				ar.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name: "Config error",
			mockConfig: func(cl *mocks.MockconfigLoader) {
				cl.EXPECT().Load(gomock.Any()).Return(nil, errors.New("test config error")).Times(1)
			},
			mockApp: func(ar *mocks.MockappRunner) {
				ar.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
			wantErr: "test config error",
		},
		{
			name: "App error",
			mockConfig: func(cl *mocks.MockconfigLoader) {
				cl.EXPECT().Load(gomock.Any()).Return(config.New(), nil).Times(1)
			},
			mockApp: func(ar *mocks.MockappRunner) {
				ar.EXPECT().Run(gomock.Any(), gomock.Any()).Return(errors.New("test app error")).Times(1)
			},
			wantErr: "test app error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConfigLoader := mocks.NewMockconfigLoader(ctrl)
			mockAppRunner := mocks.NewMockappRunner(ctrl)

			tt.mockConfig(mockConfigLoader)
			tt.mockApp(mockAppRunner)

			// arrange
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// act
			err := run(ctx, mockConfigLoader, mockAppRunner)

			// assert
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMainFunc(t *testing.T) {
	tests := []struct {
		name       string
		wantPanic  bool
		mockConfig func(loader *mocks.MockconfigLoader)
		mockApp    func(runner *mocks.MockappRunner)
	}{
		{
			name:      "Happy path",
			wantPanic: false,
			mockConfig: func(cl *mocks.MockconfigLoader) {
				cl.EXPECT().Load(gomock.Any()).Return(config.New(), nil).Times(1)
			},
			mockApp: func(ar *mocks.MockappRunner) {
				ar.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:      "Config error",
			wantPanic: true,
			mockConfig: func(cl *mocks.MockconfigLoader) {
				cl.EXPECT().Load(gomock.Any()).Return(nil, errors.New("test config error")).Times(1)
			},
			mockApp: func(ar *mocks.MockappRunner) {
				ar.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConfigLoader := mocks.NewMockconfigLoader(ctrl)
			mockAppRunner := mocks.NewMockappRunner(ctrl)

			tt.mockConfig(mockConfigLoader)
			tt.mockApp(mockAppRunner)

			cl = mockConfigLoader
			ar = mockAppRunner

			if tt.wantPanic {
				defer func() {
					r := recover()
					require.NotNil(t, r, "The code did not panic")
				}()
			}

			main()
		})
	}
}
