package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/create"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/usecase/fetch"

	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
	pb "github.com/xgmsx/go-url-shortener-ddd/proto/gen/shortener.v1"
)

type HandlerCreateLink struct {
	uc create.Usecase
}

func NewHandlerCreateLink(uc create.Usecase) *HandlerCreateLink {
	return &HandlerCreateLink{uc: uc}
}

func (h *HandlerCreateLink) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 CreateLink")
	defer span.End()

	input := dto.CreateLinkInput{URL: req.GetUrl()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.CreateLink: validate error")
		return nil, fmt.Errorf("validation error")
	}

	output, err := h.uc.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExist):
			return &pb.CreateLinkResponse{
				Url:       input.URL,
				Alias:     output.Alias,
				ExpiredAt: timestamppb.New(output.ExpiredAt),
			}, nil
		default:
			log.Error().Err(err).Msg("uc.CreateLink: internal error")
			return nil, fmt.Errorf("internal error")
		}
	}

	return &pb.CreateLinkResponse{
		Url:       input.URL,
		Alias:     output.Alias,
		ExpiredAt: timestamppb.New(output.ExpiredAt),
	}, nil
}

type HandlerFetchLink struct {
	uc fetch.Usecase
}

func NewHandlerFetchLink(uc fetch.Usecase) *HandlerFetchLink {
	return &HandlerFetchLink{uc: uc}
}

func (h *HandlerFetchLink) FetchLink(ctx context.Context, req *pb.FetchLinkRequest) (*pb.FetchLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 FetchLink")
	defer span.End()

	input := dto.FetchLinkInput{Alias: req.GetAlias()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.FetchLink: validate error")
		return nil, fmt.Errorf("validation error")
	}

	output, err := h.uc.Fetch(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.FetchLink: not found")
			return nil, fmt.Errorf("not found")
		default:
			log.Error().Err(err).Msg("uc.FetchLink: internal error")
			return nil, fmt.Errorf("internal error")
		}
	}

	return &pb.FetchLinkResponse{
		Url:       output.URL,
		Alias:     output.Alias,
		ExpiredAt: timestamppb.New(output.ExpiredAt),
	}, nil
}
