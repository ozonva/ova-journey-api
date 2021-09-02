package api

import (
	"context"
	"github.com/ozonva/ova-journey-api/internal/repo"
	"github.com/ozonva/ova-journey-api/internal/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonva/ova-journey-api/internal/models"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

// JourneyAPI - gRPC API implementation for working with journeys
type JourneyAPI struct {
	desc.UnimplementedJourneyApiV1Server
	repo      repo.Repo
	chunkSize int
}

// NewJourneyAPI returns JourneyAPI
func NewJourneyAPI(repo repo.Repo, chunkSize int) desc.JourneyApiV1Server {
	return &JourneyAPI{
		repo:      repo,
		chunkSize: chunkSize,
	}
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

	journeyID, err := a.repo.AddJourney(ctx, journey)
	if err != nil {
		log.Error().Err(err).Str("journey", journey.String()).Msg("CreateJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Str("journey", journey.String()).Msg("CreateJourneyV1: success.")
	return &desc.CreateJourneyResponseV1{JourneyId: journeyID}, nil
}

// MultiCreateJourneyV1 - create new journeys using chunks and return added journeys ids.
// If there is error for any chunk returns already added ids and error.
func (a *JourneyAPI) MultiCreateJourneyV1(ctx context.Context, req *desc.MultiCreateJourneyRequestV1) (*desc.MultiCreateJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("MultiCreateJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journeys := make([]models.Journey, len(req.Journeys))

	for i, reqJourney := range req.Journeys {
		journeys[i] = models.Journey{
			UserID:      reqJourney.UserId,
			Address:     reqJourney.Address,
			Description: reqJourney.Description,
			StartTime:   reqJourney.StartTime.AsTime(),
			EndTime:     reqJourney.EndTime.AsTime(),
		}
	}

	journeysChunks, err := utils.SplitToChunks(journeys, a.chunkSize)
	if err != nil {
		if err != nil {
			log.Error().Err(err).Msg("MultiCreateJourneyV1: failed.")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp := &desc.MultiCreateJourneyResponseV1{}
	for _, chunk := range journeysChunks {
		ids, err := a.repo.MultiAddJourneys(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msg("MultiCreateJourneyV1: failed.")
			return resp, status.Error(codes.Internal, err.Error())
		}
		resp.JourneyIds = append(resp.JourneyIds, ids...)
	}

	log.Debug().Msg("MultiCreateJourneyV1: success.")

	return resp, nil
}

// DescribeJourneyV1 - get journey description by journeyID
func (a *JourneyAPI) DescribeJourneyV1(ctx context.Context, req *desc.DescribeJourneyRequestV1) (*desc.DescribeJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journey, err := a.repo.DescribeJourney(ctx, req.JourneyId)
	if err != nil {
		log.Error().Err(err).Uint64("journeyId", req.JourneyId).Msg("DescribeJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.JourneyId).Msg("DescribeJourneyV1: success.")
	return &desc.DescribeJourneyResponseV1{
		Journey: &desc.Journey{
			JourneyId:   journey.JourneyID,
			UserId:      journey.UserID,
			Address:     journey.Address,
			Description: journey.Description,
			StartTime:   timestamppb.New(journey.StartTime),
			EndTime:     timestamppb.New(journey.EndTime),
		},
	}, nil
}

// ListJourneysV1 - get list of journey with offset and limit
func (a *JourneyAPI) ListJourneysV1(ctx context.Context, req *desc.ListJourneysRequestV1) (*desc.ListJourneysResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListJourneysV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journeys, err := a.repo.ListJourneys(ctx, req.Limit, req.Offset)
	if err != nil {
		log.Error().Err(err).Uint64("offset", req.Offset).Uint64("limit", req.Limit).Msg("ListJourneysV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &desc.ListJourneysResponseV1{Journeys: make([]*desc.Journey, len(journeys))}
	for i, journey := range journeys {
		resp.Journeys[i] = &desc.Journey{
			JourneyId:   journey.JourneyID,
			UserId:      journey.UserID,
			Address:     journey.Address,
			Description: journey.Description,
			StartTime:   timestamppb.New(journey.StartTime),
			EndTime:     timestamppb.New(journey.EndTime),
		}
	}

	log.Debug().Uint64("offset", req.Offset).Uint64("limit", req.Limit).Msg("ListJourneysV1: success.")
	return resp, nil
}

// RemoveJourneyV1 - remove journey
func (a *JourneyAPI) RemoveJourneyV1(ctx context.Context, req *desc.RemoveJourneyRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := a.repo.RemoveJourney(ctx, req.JourneyId); err != nil {
		log.Error().Err(err).Uint64("journeyId", req.JourneyId).Msg("RemoveJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.JourneyId).Msg("RemoveJourneyV1: success.")
	return &emptypb.Empty{}, nil
}

// UpdateJourneyV1 - find journey by id and update another fields
func (a *JourneyAPI) UpdateJourneyV1(ctx context.Context, req *desc.UpdateJourneyRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("UpdateJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journey := models.Journey{
		JourneyID:   req.Journey.JourneyId,
		UserID:      req.Journey.UserId,
		Address:     req.Journey.Address,
		Description: req.Journey.Description,
		StartTime:   req.Journey.StartTime.AsTime(),
		EndTime:     req.Journey.EndTime.AsTime(),
	}
	if err := a.repo.UpdateJourney(ctx, journey); err != nil {
		log.Error().Err(err).Str("journey", req.Journey.String()).Msg("UpdateJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.Journey.JourneyId).Msg("UpdateJourneyV1: success.")
	return &emptypb.Empty{}, nil
}
