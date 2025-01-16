package grpc_controller //nolint:stylecheck

import (
	"google.golang.org/grpc"

	pb "github.com/xgmsx/go-url-shortener-ddd/generated/protobuf/shortener.v1"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/usecase"
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
