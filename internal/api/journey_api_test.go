package api

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-journey-api/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/ozonva/ova-journey-api/internal/mocks"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

var _ = Describe("JourneyApi", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockRepo
		api      desc.JourneyApiV1Server
		ctx      context.Context

		timeStart     = time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)
		timeEnd       = time.Date(2021, 01, 02, 0, 0, 0, 0, time.UTC)
		journeysTable = []models.Journey{
			{JourneyID: 0, UserID: 1, Address: "Воронеж", Description: "", StartTime: timeStart, EndTime: timeEnd},
			{JourneyID: 1, UserID: 1, Address: "Уфа", Description: "", StartTime: timeStart, EndTime: timeEnd},
			{JourneyID: 2, UserID: 2, Address: "Москва", Description: "", StartTime: timeStart, EndTime: timeEnd},
		}
		errRepo = errors.New("repo error")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
		ctx = context.Background()
	})

	JustBeforeEach(func() {
		api = NewJourneyAPI(mockRepo)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("CreateJourneyV1", func() {
		Context("Success create journey", func() {
			It("should return success result with new journey id", func() {
				newJourneyID := uint64(1)
				mockRepo.EXPECT().AddJourney(ctx, journeysTable[0]).Return(newJourneyID, nil).Times(1)

				result, err := api.CreateJourneyV1(ctx, &desc.CreateJourneyRequestV1{
					UserId:    journeysTable[0].UserID,
					Address:   journeysTable[0].Address,
					StartTime: timestamppb.New(journeysTable[0].StartTime),
					EndTime:   timestamppb.New(journeysTable[0].EndTime),
				})

				Expect(result.JourneyId).Should(Equal(newJourneyID))
				Expect(err).Should(BeNil())
			})
		})

		Context("Incorrect journey in request", func() {
			It("should return error without calling repo", func() {
				mockRepo.EXPECT().AddJourney(ctx, gomock.Any()).Times(0)

				result, err := api.CreateJourneyV1(ctx, &desc.CreateJourneyRequestV1{
					UserId:    0,
					Address:   journeysTable[0].Address,
					StartTime: timestamppb.New(journeysTable[0].StartTime),
					EndTime:   timestamppb.New(journeysTable[0].EndTime),
				})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("Error in repo", func() {
			It("should return error", func() {
				mockRepo.EXPECT().AddJourney(ctx, journeysTable[0]).Return(uint64(0), errRepo).Times(1)

				result, err := api.CreateJourneyV1(ctx, &desc.CreateJourneyRequestV1{
					UserId:    journeysTable[0].UserID,
					Address:   journeysTable[0].Address,
					StartTime: timestamppb.New(journeysTable[0].StartTime),
					EndTime:   timestamppb.New(journeysTable[0].EndTime),
				})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("DescribeJourneyV1", func() {
		Context("Success describe journey", func() {
			It("should return success result with journey", func() {
				journeyID := uint64(1)
				mockRepo.EXPECT().DescribeJourney(ctx, journeyID).Return(&journeysTable[1], nil).Times(1)

				result, err := api.DescribeJourneyV1(ctx, &desc.DescribeJourneyRequestV1{JourneyId: journeyID})

				Expect(result.Journey.JourneyId).Should(Equal(journeysTable[1].JourneyID))
				Expect(result.Journey.UserId).Should(Equal(journeysTable[1].UserID))
				Expect(result.Journey.Address).Should(Equal(journeysTable[1].Address))
				Expect(result.Journey.Description).Should(Equal(journeysTable[1].Description))
				Expect(result.Journey.StartTime.AsTime()).Should(Equal(journeysTable[1].StartTime))
				Expect(result.Journey.EndTime.AsTime()).Should(Equal(journeysTable[1].EndTime))
				Expect(err).Should(BeNil())
			})
		})

		Context("Incorrect journeyId in request", func() {
			It("should return error without calling repo", func() {
				journeyID := uint64(0)
				mockRepo.EXPECT().DescribeJourney(ctx, gomock.Any()).Times(0)

				result, err := api.DescribeJourneyV1(ctx, &desc.DescribeJourneyRequestV1{JourneyId: journeyID})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("Error in repo", func() {
			It("should return error", func() {
				journeyID := uint64(1)
				mockRepo.EXPECT().DescribeJourney(ctx, journeyID).Return(nil, errRepo).Times(1)

				result, err := api.DescribeJourneyV1(ctx, &desc.DescribeJourneyRequestV1{JourneyId: journeyID})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("ListJourneysV1", func() {
		Context("Success get list of journeys", func() {
			It("should return success result with journeys", func() {
				limit := uint64(2)
				offset := uint64(1)
				mockRepo.EXPECT().ListJourneys(ctx, limit, offset).Return(journeysTable[1:3], nil).Times(1)

				result, err := api.ListJourneysV1(ctx, &desc.ListJourneysRequestV1{Offset: offset, Limit: limit})

				Expect(result.Journeys).Should(HaveLen(2))
				Expect(err).Should(BeNil())
			})
		})

		Context("Incorrect limit in request", func() {
			It("should return error without calling repo", func() {
				limit := uint64(0)
				offset := uint64(1)
				mockRepo.EXPECT().ListJourneys(ctx, gomock.Any(), gomock.Any()).Times(0)

				result, err := api.ListJourneysV1(ctx, &desc.ListJourneysRequestV1{Offset: offset, Limit: limit})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("Error in repo", func() {
			It("should return error", func() {
				limit := uint64(1)
				offset := uint64(1)
				mockRepo.EXPECT().ListJourneys(ctx, limit, offset).Return(nil, errRepo).Times(1)

				result, err := api.ListJourneysV1(ctx, &desc.ListJourneysRequestV1{Offset: offset, Limit: limit})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("RemoveJourneyV1", func() {
		Context("Success remove journey", func() {
			It("should return success empty result", func() {
				journeyID := uint64(1)

				mockRepo.EXPECT().RemoveJourney(ctx, journeyID).Return(nil).Times(1)

				result, err := api.RemoveJourneyV1(ctx, &desc.RemoveJourneyRequestV1{JourneyId: journeyID})

				Expect(result).Should(Equal(&emptypb.Empty{}))
				Expect(err).Should(BeNil())
			})
		})

		Context("Incorrect journeyId in request", func() {
			It("should return error without calling repo", func() {
				journeyID := uint64(0)
				mockRepo.EXPECT().RemoveJourney(ctx, gomock.Any()).Times(0)

				result, err := api.RemoveJourneyV1(ctx, &desc.RemoveJourneyRequestV1{JourneyId: journeyID})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("Error in repo", func() {
			It("should return error", func() {
				mockRepo.EXPECT().RemoveJourney(ctx, gomock.Any()).Return(errRepo).Times(1)

				result, err := api.RemoveJourneyV1(ctx, &desc.RemoveJourneyRequestV1{JourneyId: 1})

				Expect(result).Should(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
