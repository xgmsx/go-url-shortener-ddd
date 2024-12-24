package grpc

import (
	"net"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Registrable interface {
	Register(s *grpc.Server)
}

type Config struct {
	Port string `env:"GRPC_PORT, default=50051"`
}

type Server struct {
	srv    *grpc.Server
	config Config
	notify chan error
}

func New(ch chan error, config Config, controllers ...Registrable) *Server {
	s := &Server{
		srv:    grpc.NewServer(),
		config: config,
		notify: ch,
	}

	for _, controller := range controllers {
		controller.Register(s.srv)
	}
	reflection.Register(s.srv)

	go func() {
		lis, err := net.Listen("tcp", ":"+config.Port)
		if err != nil {
			s.notify <- err
			return
		}
		s.notify <- s.srv.Serve(lis)
	}()

	log.Info().Msg("gRPC server started on port: " + config.Port)

	return s
}

func (s *Server) Close() {
	s.srv.GracefulStop()
	log.Info().Msg("gRPC server closed")
}

func (s *Server) Notify(err error) {
	if err != nil {
		s.notify <- err
	}
}
