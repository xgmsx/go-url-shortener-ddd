package app

import (
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/grpc_controller"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/http_router"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/kafka_consumer"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/grpc"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/kafka/reader"
)

func getHTTPController(ch chan error, c *config.Config, uc *usecase.UseCase) *http.Server {
	return http.New(ch, c.HTTP, http.BuildOptions(c.HTTP,
		http.WithMiddleware(requestid.New(requestid.Config{}), c.HTTP.UseRequestID),
		http.WithMiddleware(pprof.New(pprof.Config{}), c.HTTP.UsePprof),
		http.WithMiddleware(recover.New(recover.Config{}), c.HTTP.UseRecover),
		http.WithRouter("/", http.DefaultRouter{}),
		http.WithRouter("/api/shortener", http_router.New(uc)),
	))
}

func getGRPCController(ch chan error, c *config.Config, uc *usecase.UseCase) *grpc.Server {
	return grpc.New(ch, c.GRPC, grpc_controller.New(uc))
}

func getKafkaController(ch chan error, r *reader.Reader, uc *usecase.UseCase) *kafka_consumer.Consumer {
	return kafka_consumer.New(ch, r, uc)
}
