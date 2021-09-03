package api

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-journey-api/internal/kafka"
	"github.com/ozonva/ova-journey-api/internal/metrics"
	"github.com/ozonva/ova-journey-api/internal/models"
	"github.com/ozonva/ova-journey-api/internal/repo"
	"github.com/ozonva/ova-journey-api/internal/utils"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// JourneyAPI - gRPC API implementation for working with journeys
type JourneyAPI struct {
	desc.UnimplementedJourneyApiV1Server
	repo      repo.Repo
	producer  kafka.Producer
	metric    metrics.Metrics
	chunkSize int
}

// NewJourneyAPI returns JourneyAPI
func NewJourneyAPI(repo repo.Repo, producer kafka.Producer, metric metrics.Metrics, chunkSize int) desc.JourneyApiV1Server {
	return &JourneyAPI{
		repo:      repo,
		producer:  producer,
		chunkSize: chunkSize,
		metric:    metric,
	}
}

// CreateJourneyV1 - create new journey
func (api *JourneyAPI) CreateJourneyV1(ctx context.Context, req *desc.CreateJourneyRequestV1) (*desc.CreateJourneyResponseV1, error) {
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

	journeyID, err := api.repo.AddJourney(ctx, journey)
	if err != nil {
		log.Error().Err(err).Str("journey", journey.String()).Msg("CreateJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Str("journey", journey.String()).Msg("CreateJourneyV1: success.")
	api.metric.CreateJourneyCounterInc()

	return &desc.CreateJourneyResponseV1{JourneyId: journeyID}, nil
}

// MultiCreateJourneyV1 - create new journeys using chunks and return added journeys ids.
// If there is error for any chunk returns already added ids and error.
func (api *JourneyAPI) MultiCreateJourneyV1(ctx context.Context, req *desc.MultiCreateJourneyRequestV1) (*desc.MultiCreateJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("MultiCreateJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreateJourneyV1")
	defer span.Finish()

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

	journeysChunks, err := utils.SplitToChunks(journeys, api.chunkSize)
	if err != nil {
		log.Error().Err(err).Msg("MultiCreateJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &desc.MultiCreateJourneyResponseV1{}
	for _, chunk := range journeysChunks {
		ids, err := api.repo.MultiAddJourneys(ctx, chunk)
		if err != nil {
			log.Error().Err(err).Msg("MultiCreateJourneyV1: failed.")
			return resp, status.Error(codes.Internal, err.Error())
		}
		resp.JourneyIds = append(resp.JourneyIds, ids...)

		childSpan := tracer.StartSpan(
			"MultiCreateJourneyV1: chunk",
			opentracing.Tag{Key: "chunkSize", Value: len(chunk)},
			opentracing.ChildOf(span.Context()),
		)
		childSpan.Finish()
	}

	log.Debug().Msg("MultiCreateJourneyV1: success.")
	api.metric.MultiCreateJourneyCounterInc()

	return resp, nil
}

// DescribeJourneyV1 - get journey description by journeyID
func (api *JourneyAPI) DescribeJourneyV1(ctx context.Context, req *desc.DescribeJourneyRequestV1) (*desc.DescribeJourneyResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journey, err := api.repo.DescribeJourney(ctx, req.JourneyId)
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
func (api *JourneyAPI) ListJourneysV1(ctx context.Context, req *desc.ListJourneysRequestV1) (*desc.ListJourneysResponseV1, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListJourneysV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journeys, err := api.repo.ListJourneys(ctx, req.Limit, req.Offset)
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
func (api *JourneyAPI) RemoveJourneyV1(ctx context.Context, req *desc.RemoveJourneyRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveJourneyV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := api.repo.RemoveJourney(ctx, req.JourneyId); err != nil {
		log.Error().Err(err).Uint64("journeyId", req.JourneyId).Msg("RemoveJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.JourneyId).Msg("RemoveJourneyV1: success.")
	api.metric.DeleteJourneyCounterInc()

	return &emptypb.Empty{}, nil
}

// UpdateJourneyV1 - find journey by id and update another fields
func (api *JourneyAPI) UpdateJourneyV1(ctx context.Context, req *desc.UpdateJourneyRequestV1) (*emptypb.Empty, error) {
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
	if err := api.repo.UpdateJourney(ctx, journey); err != nil {
		log.Error().Err(err).Str("journey", req.Journey.String()).Msg("UpdateJourneyV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.Journey.JourneyId).Msg("UpdateJourneyV1: success.")
	api.metric.UpdateJourneyCounterInc()

	return &emptypb.Empty{}, nil
}

// CreateJourneyTaskV1 - create new journey using producer
func (api *JourneyAPI) CreateJourneyTaskV1(ctx context.Context, req *desc.CreateJourneyTaskRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateJourneyTaskV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	journey := models.Journey{
		UserID:      req.UserId,
		Address:     req.Address,
		Description: req.Description,
		StartTime:   req.StartTime.AsTime(),
		EndTime:     req.EndTime.AsTime(),
	}

	err := api.producer.Send(kafka.Message{
		MessageType: kafka.CreateJourney,
		Value:       journey,
	})

	if err != nil {
		log.Error().Err(err).Str("journey", journey.String()).Msg("CreateJourneyTaskV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Str("journey", journey.String()).Msg("CreateJourneyTaskV1: success.")
	api.metric.CreateJourneyCounterInc()

	return &emptypb.Empty{}, nil
}

// MultiCreateJourneyTaskV1 - create new journeys using producer and splitting on chunks
func (api *JourneyAPI) MultiCreateJourneyTaskV1(ctx context.Context, req *desc.MultiCreateJourneyTaskRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("MultiCreateJourneyTaskV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreateJourneyTaskV1")
	defer span.Finish()

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

	journeysChunks, err := utils.SplitToChunks(journeys, api.chunkSize)
	if err != nil {
		if err != nil {
			log.Error().Err(err).Msg("MultiCreateJourneyTaskV1: failed.")
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, chunk := range journeysChunks {
		err = api.producer.Send(kafka.Message{
			MessageType: kafka.MultiCreateJourney,
			Value:       chunk,
		})
		if err != nil {
			log.Error().Err(err).Msg("MultiCreateJourneyTaskV1: failed.")
			return nil, status.Error(codes.Internal, err.Error())
		}

		childSpan := tracer.StartSpan(
			"MultiCreateJourneyTaskV1: chunk",
			opentracing.Tag{Key: "chunkSize", Value: len(chunk)},
			opentracing.ChildOf(span.Context()),
		)
		childSpan.Finish()
	}

	log.Debug().Msg("MultiCreateJourneyTaskV1: success send to producer.")
	api.metric.MultiCreateJourneyCounterInc()

	return &emptypb.Empty{}, nil
}

// RemoveJourneyTaskV1 - remove journey using producer
func (api *JourneyAPI) RemoveJourneyTaskV1(ctx context.Context, req *desc.RemoveJourneyTaskRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveJourneyTaskV1: invalid request.")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := api.producer.Send(kafka.Message{
		MessageType: kafka.DeleteJourney,
		Value:       req.JourneyId,
	})
	if err != nil {
		log.Error().Err(err).Uint64("journeyId", req.JourneyId).Msg("RemoveJourneyTaskV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("RemoveJourneyTaskV1: success.")
	api.metric.DeleteJourneyCounterInc()

	return &emptypb.Empty{}, nil
}

// UpdateJourneyTaskV1 - find journey by id and update another fields using producer
func (api *JourneyAPI) UpdateJourneyTaskV1(ctx context.Context, req *desc.UpdateJourneyTaskRequestV1) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("UpdateJourneyTaskV1: invalid request.")
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
	err := api.producer.Send(kafka.Message{
		MessageType: kafka.UpdateJourney,
		Value:       journey,
	})
	if err != nil {
		log.Error().Err(err).Str("journey", req.Journey.String()).Msg("UpdateJourneyTaskV1: failed.")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Uint64("journeyId", req.Journey.JourneyId).Msg("UpdateJourneyTaskV1: success.")
	api.metric.UpdateJourneyCounterInc()

	return &emptypb.Empty{}, nil
}
