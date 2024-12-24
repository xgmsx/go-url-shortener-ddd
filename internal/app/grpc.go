package app

import (
	"github.com/xgmsx/go-url-shortener-ddd/internal/config"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/controller/grpc_controller"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/grpc"
)

func getGRPCController(ch chan error, c *config.Config, uc *usecase.UseCase) *grpc.Server {
	return grpc.New(ch, c.GRPC, grpc_controller.New(uc))
}
