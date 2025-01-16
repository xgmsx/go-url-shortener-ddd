package grpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"
	pb "github.com/xgmsx/go-url-shortener-ddd/proto/gen/shortener.v1"
)

type Controller struct {
	pb.UnimplementedShortenerServer
	createHandler *HandlerCreateLink
	fetchHandler  *HandlerFetchLink
}

func New(ucCreate create.Usecase, ucFetch fetch.Usecase) *Controller {
	return &Controller{
		createHandler: NewHandlerCreateLink(ucCreate),
		fetchHandler:  NewHandlerFetchLink(ucFetch),
	}
}

func (c *Controller) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	return c.createHandler.CreateLink(ctx, req)
}

func (c *Controller) FetchLink(ctx context.Context, req *pb.FetchLinkRequest) (*pb.FetchLinkResponse, error) {
	return c.fetchHandler.FetchLink(ctx, req)
}

func (c *Controller) Register(server *grpc.Server) {
	pb.RegisterShortenerServer(server, c)
}
