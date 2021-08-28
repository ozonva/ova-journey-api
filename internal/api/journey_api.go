package api

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ozonva/ova-journey-api/internal/models"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

const msgUnimplemented = "method is not implemented"

// JourneyAPI - gRPC API implementation for working with journeys
type JourneyAPI struct {
	desc.UnimplementedJourneyApiV1Server
}

// NewJourneyAPI returns JourneyAPI
func NewJourneyAPI() desc.JourneyApiV1Server {
	return &JourneyAPI{}
}

// CreateJourneyV1 - create new journey
func (a *JourneyAPI) CreateJourneyV1(ctx context.Context, req *desc.CreateJourneyRequestV1) (*desc.CreateJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journey := models.Journey{
		UserID:      req.UserId,
		Address:     req.Address,
		Description: req.Description,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
	}

	log.Debug().Str("journey", journey.String()).Msg("CreateJourneyV1: success.")

	return nil, status.Error(codes.Unimplemented, msgUnimplemented)
}

// DescribeJourneyV1 - get journey description by journeyID
func (a *JourneyAPI) DescribeJourneyV1(ctx context.Context, req *desc.DescribeJourneyRequestV1) (*desc.DescribeJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Debug().Uint64("journeyId", req.JourneyId).Msg("DescribeJourneyV1: success.")

	return nil, status.Error(codes.Unimplemented, msgUnimplemented)
}

// ListJourneysV1 - get list of journey with offset and limit
func (a *JourneyAPI) ListJourneysV1(ctx context.Context, req *desc.ListJourneysRequestV1) (*desc.ListJourneysResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListJourneysV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Debug().Uint64("offset", req.Offset).Uint64("limit", req.Limit).Msg("ListJourneysV1: success.")

	return nil, status.Error(codes.Unimplemented, msgUnimplemented)
}

// RemoveJourneyV1 - remove journey
func (a *JourneyAPI) RemoveJourneyV1(ctx context.Context, req *desc.RemoveJourneyRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Debug().Uint64("offset", req.JourneyId).Msg("RemoveJourneyV1: success.")

	return nil, status.Error(codes.Unimplemented, msgUnimplemented)
}
