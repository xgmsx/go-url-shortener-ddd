package grpc

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type registrable interface {
	Register(s *grpc.Server)
}

type Config struct {
	Port string `env:"GRPC_PORT, default=50051"`
}

type Server struct {
	srv *grpc.Server
}

func New(controllers ...registrable) *Server {
	srv := grpc.NewServer()

	for _, c := range controllers {
		c.Register(srv)
	}
	reflection.Register(srv)

	return &Server{srv: srv}
}

func (s *Server) Serve(ctx context.Context, port string) error {
	var lc net.ListenConfig

	lis, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Info().Msg("gRPC server started on port: " + port)
	return s.srv.Serve(lis)
}

func (s *Server) Close() {
	s.srv.GracefulStop()
	log.Info().Msg("gRPC server closed")
}
