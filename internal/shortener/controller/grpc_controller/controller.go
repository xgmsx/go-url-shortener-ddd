package grpc_controller

import (
	"google.golang.org/grpc"

	pb "url-shortener/generated/protobuf/shortener.v1"
	"url-shortener/internal/shortener/usecase"
)

type Controller struct {
	pb.UnimplementedShortenerServer
	uc *usecase.UseCase
}

func New(uc *usecase.UseCase) Controller {
	return Controller{uc: uc}
}

func (c Controller) Register(server *grpc.Server) {
	pb.RegisterShortenerServer(server, c)
}
