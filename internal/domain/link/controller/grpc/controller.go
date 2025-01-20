package grpc

import (
	"google.golang.org/grpc"

	pb "github.com/xgmsx/go-url-shortener-ddd/generated/protobuf/shortener.v1"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/usecase/fetch"
)

type Controller struct {
	pb.UnimplementedShortenerServer
	ucCreate create.Usecase
	ucFetch  fetch.Usecase
}

func New(ucCreate create.Usecase, ucFetch fetch.Usecase) *Controller {
	return &Controller{
		ucCreate: ucCreate,
		ucFetch:  ucFetch,
	}
}

func (c *Controller) Register(server *grpc.Server) {
	pb.RegisterShortenerServer(server, c)
}
