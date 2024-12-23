package app

import (
	"url-shortener/internal/config"
	"url-shortener/internal/shortener/controller/http_router"
	"url-shortener/internal/shortener/usecase"
	"url-shortener/pkg/http"

	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func getHTTPController(ch chan error, c config.Config, uc *usecase.UseCase) *http.Server {
	return http.New(ch, c.HTTP, http.BuildOptions(c.HTTP,
		http.WithMiddleware(requestid.New(requestid.Config{})),
		http.WithMiddleware(pprof.New(pprof.Config{Prefix: "/pprof"})),
		http.WithMiddleware(recover.New(recover.Config{})),
		http.WithRouter("/", http.DefaultRouter{}),
		http.WithRouter("/api/shortener", http_router.New(uc)),
	))
}
