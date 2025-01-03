package grpc_controller //nolint:stylecheck

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/xgmsx/go-url-shortener-ddd/generated/protobuf/shortener.v1"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/dto"
	"github.com/xgmsx/go-url-shortener-ddd/internal/shortener/entity"
	"github.com/xgmsx/go-url-shortener-ddd/pkg/observability/otel/tracer"
)

func (c Controller) CreateLink(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 CreateLink")
	defer span.End()

	input := dto.CreateLinkInput{URL: req.GetUrl()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.CreateLink: validate error")
		return nil, err
	}

	output, err := c.uc.CreateLink(ctx, input)
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

func (c Controller) GetLink(ctx context.Context, req *pb.GetLinkRequest) (*pb.GetLinkResponse, error) {
	ctx, span := tracer.Start(ctx, "grpc/v1 GetLink")
	defer span.End()

	input := dto.GetLinkInput{Alias: req.GetAlias()}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Msg("uc.GetLink: validate error")
		return nil, err
	}

	output, err := c.uc.GetLink(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			log.Error().Err(err).Msg("uc.GetLink: not found")
			return nil, fmt.Errorf("not found")
		default:
			log.Error().Err(err).Msg("uc.GetLink: internal error")
			return nil, fmt.Errorf("internal error")
		}
	}

	return &pb.GetLinkResponse{
		Url:       output.URL,
		Alias:     output.Alias,
		ExpiredAt: timestamppb.New(output.ExpiredAt),
	}, nil
}
