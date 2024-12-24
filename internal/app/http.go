package app

import (
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/http_router"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/http"
)

func getHTTPController(ch chan error, c *config.Config, uc *usecase.UseCase) *http.Server {
	return http.New(ch, c.HTTP, http.BuildOptions(c.HTTP,
		http.WithMiddleware(requestid.New(requestid.Config{})),
		http.WithMiddleware(pprof.New(pprof.Config{Prefix: "/pprof"})),
		http.WithMiddleware(recover.New(recover.Config{})),
		http.WithRouter("/", http.DefaultRouter{}),
		http.WithRouter("/api/shortener", http_router.New(uc)),
	))
}
