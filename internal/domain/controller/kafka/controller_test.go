package kafka_test

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/mock/gomock"

	controllerKafka "github.com/xgmsx/go-url-shortener-ddd/internal/domain/controller/kafka"
	mocksReader "github.com/xgmsx/go-url-shortener-ddd/internal/domain/controller/kafka/mocks"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	ucCreate "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	mocksCreate "github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create/mocks"
)

func TestKafkaController(t *testing.T) {
	initMock := func() (*gomock.Controller, *mocksCreate.Mockdatabase, *mocksCreate.Mockcache, *mocksCreate.Mockpublisher) {
		ctrl := gomock.NewController(t)
		publisher := mocksCreate.NewMockpublisher(ctrl)
		database := mocksCreate.NewMockdatabase(ctrl)
		cache := mocksCreate.NewMockcache(ctrl)
		return ctrl, database, cache, publisher
	}

	testCases := []struct {
		name      string
		input     string
		setupMock func(*gomock.Controller, *mocksCreate.Mockdatabase, *mocksCreate.Mockcache, *mocksCreate.Mockpublisher)
	}{
		{
			name:  "Happy path",
			input: `{"url": "https://example.com"}`,
			setupMock: func(ctrl *gomock.Controller, database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				cache.EXPECT().PutLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				publisher.EXPECT().SendLink(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
		},
		{
			name:  "Link already exists",
			input: `{"url": "https://example.com"}`,
			setupMock: func(ctrl *gomock.Controller, database *mocksCreate.Mockdatabase, cache *mocksCreate.Mockcache, publisher *mocksCreate.Mockpublisher) {
				database.EXPECT().CreateLink(gomock.Any(), gomock.Any()).Return(entity.ErrAlreadyExist).Times(1)
			},
		},
		{
			name:  "Validation error",
			input: `https://example.com`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			ctrl, database, cache, publisher := initMock()
			defer ctrl.Finish()

			if tc.setupMock != nil {
				tc.setupMock(ctrl, database, cache, publisher)
			}

			// arrange
			msg := kafka.Message{Value: []byte(tc.input)}
			doFunc := func(ctx context.Context) (kafka.Message, error) { cancel(); return msg, nil }
			reader := mocksReader.NewMockkafkaReader(ctrl)
			reader.EXPECT().FetchMessage(gomock.Any()).DoAndReturn(doFunc).Times(1)
			reader.EXPECT().CommitMessages(ctx, msg).Return(nil).MaxTimes(1)

			// act
			controller := controllerKafka.New(reader, ucCreate.New(database, cache, publisher))
			go controller.Consume(ctx)

			<-time.After(time.Millisecond * 50)
		})
	}
}
