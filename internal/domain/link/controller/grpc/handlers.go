package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/xgmsx/go-url-shortener-ddd/generated/protobuf/shortener.v1"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/domain/link/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

func (c *Controller) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 CreateLink")
	defer span.End()

	input := dto.CreateLinkInput{URL: req.GetUrl()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.CreateLink: validate error")
		return nil, err
	}

	output, err := c.ucCreate.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrAlreadyExist):
			return &pb.CreateLinkResponse{
				Url:       input.URL,
				Alias:     output.Alias,
				ExpiredAt: timestamppb.New(output.ExpiredAt),
			}, nil
		case errors.Is(err, entity.ErrEntityValidation):
			log.Error().Err(err).Msg("uc.CreateLink: validate error")
			return nil, fmt.Errorf("validate error")
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

func (c *Controller) FetchLink(ctx context.Context, req *pb.FetchLinkRequest) (*pb.FetchLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 FetchLink")
	defer span.End()

	input := dto.FetchLinkInput{Alias: req.GetAlias()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.FetchLink: validate error")
		return nil, err
	}

	output, err := c.ucFetch.Fetch(ctx, input)
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
