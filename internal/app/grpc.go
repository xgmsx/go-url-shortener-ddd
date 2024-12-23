package app

import (
	"url-shortener/internal/config"
	"url-shortener/internal/shortener/controller/grpc_controller"
	"url-shortener/internal/shortener/usecase"
	"url-shortener/pkg/grpc"
)

func getGRPCController(ch chan error, c config.Config, uc *usecase.UseCase) *grpc.Server {
	return grpc.New(ch, c.GRPC, grpc_controller.New(uc))
}
