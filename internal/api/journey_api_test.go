package api

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-journey-api/internal/kafka"
	"github.com/ozonva/ova-journey-api/internal/models"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/ozonva/ova-journey-api/internal/mocks"
	desc "github.com/ozonva/ova-journey-api/pkg/ova-journey-api"
)

var _ = Describe("JourneyApi", func() {
	var (
		ctrl         *gomock.Controller
		mockRepo     *mocks.MockRepo
		mockProducer *mocks.MockProducer
		mockMetrics  *mocks.MockMetrics
		api          desc.JourneyApiV1Server
		ctx          context.Context

		chunkSize     = 2
		timeStart     = time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)
		timeEnd       = time.Date(2021, 01, 02, 0, 0, 0, 0, time.UTC)
		journeysTable = []models.Journey{
			{JourneyID: 0, UserID: 1, Address: "Воронеж", Description: "", StartTime: timeStart, EndTime: timeEnd},
			{JourneyID: 1, UserID: 1, Address: "Уфа", Description: "", StartTime: timeStart, EndTime: timeEnd},
			{JourneyID: 2, UserID: 2, Address: "Москва", Description: "", StartTime: timeStart, EndTime: timeEnd},
		}
		errRepo     = errors.New("repo error")
		errProducer = errors.New("producer error")
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
		mockProducer = mocks.NewMockProducer(ctrl)
		mockMetrics = mocks.NewMockMetrics(ctrl)
		ctx = context.Background()
	})

	JustBeforeEach(func() {
		api = NewJourneyAPI(mockRepo, mockProducer, mockMetrics, chunkSize)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Using repo", func() {
		Context("CreateJourneyV1", func() {
			Context("Success create journey", func() {
				It("should return success result with new journey id", func() {
					newJourneyID := uint64(1)
					mockRepo.EXPECT().AddJourney(ctx, journeysTable[0]).Return(newJourneyID, nil).Times(1)
					mockMetrics.EXPECT().CreateJourneyCounterInc().Times(1)

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

		Context("MultiCreateJourneyV1", func() {
			var journeysTableZeroIds []models.Journey

			BeforeEach(func() {
				journeysTableZeroIds = make([]models.Journey, len(journeysTable))
				copy(journeysTableZeroIds, journeysTable)
				for i := range journeysTableZeroIds {
					journeysTableZeroIds[i].JourneyID = 0
				}
			})

			Context("Success create journeys", func() {
				It("should return success result with new journey ids", func() {
					newJourneyIDs := []uint64{1, 2, 3}

					mockRepo.EXPECT().MultiAddJourneys(ctx, journeysTableZeroIds[0:2]).Return(newJourneyIDs[0:2], nil).Times(1)
					mockRepo.EXPECT().MultiAddJourneys(ctx, journeysTableZeroIds[2:]).Return(newJourneyIDs[2:], nil).Times(1)
					mockMetrics.EXPECT().MultiCreateJourneyCounterInc().Times(1)

					req := &desc.MultiCreateJourneyRequestV1{}
					for _, journey := range journeysTable {
						req.Journeys = append(req.Journeys, &desc.CreateJourneyRequestV1{
							UserId:    journey.UserID,
							Address:   journey.Address,
							StartTime: timestamppb.New(journey.StartTime),
							EndTime:   timestamppb.New(journey.EndTime),
						})
					}

					result, err := api.MultiCreateJourneyV1(ctx, req)

					Expect(result.JourneyIds).Should(Equal(newJourneyIDs))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journeys count in request", func() {
				It("should return error without calling repo", func() {
					mockRepo.EXPECT().MultiAddJourneys(ctx, gomock.Any()).Times(0)

					result, err := api.MultiCreateJourneyV1(ctx, &desc.MultiCreateJourneyRequestV1{})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in repo", func() {
				It("should return error with added ids", func() {
					newJourneyIDs := []uint64{1, 2, 3}

					mockRepo.EXPECT().MultiAddJourneys(ctx, journeysTableZeroIds[0:2]).Return(newJourneyIDs[0:2], nil).Times(1)
					mockRepo.EXPECT().MultiAddJourneys(ctx, journeysTableZeroIds[2:]).Return(nil, errRepo).Times(1)

					req := &desc.MultiCreateJourneyRequestV1{}
					for _, journey := range journeysTable {
						req.Journeys = append(req.Journeys, &desc.CreateJourneyRequestV1{
							UserId:    journey.UserID,
							Address:   journey.Address,
							StartTime: timestamppb.New(journey.StartTime),
							EndTime:   timestamppb.New(journey.EndTime),
						})
					}

					result, err := api.MultiCreateJourneyV1(ctx, req)

					Expect(result.JourneyIds).Should(Equal(newJourneyIDs[0:2]))
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
					mockMetrics.EXPECT().DeleteJourneyCounterInc().Times(1)

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

		Context("UpdateJourneyV1", func() {
			Context("Success update journey", func() {
				It("should return success empty result", func() {
					mockRepo.EXPECT().UpdateJourney(ctx, journeysTable[2]).Return(nil).Times(1)
					mockMetrics.EXPECT().UpdateJourneyCounterInc().Times(1)

					result, err := api.UpdateJourneyV1(ctx, &desc.UpdateJourneyRequestV1{
						Journey: &desc.Journey{
							JourneyId: journeysTable[2].JourneyID,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(Equal(&emptypb.Empty{}))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journey id in request", func() {
				It("should return error without calling repo", func() {
					mockRepo.EXPECT().UpdateJourney(ctx, gomock.Any()).Times(0)

					result, err := api.UpdateJourneyV1(ctx, &desc.UpdateJourneyRequestV1{
						Journey: &desc.Journey{
							JourneyId: 0,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in repo", func() {
				It("should return error", func() {
					mockRepo.EXPECT().UpdateJourney(ctx, journeysTable[2]).Return(errRepo).Times(1)

					result, err := api.UpdateJourneyV1(ctx, &desc.UpdateJourneyRequestV1{
						Journey: &desc.Journey{
							JourneyId: journeysTable[2].JourneyID,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

	Context("Using producer", func() {
		Context("CreateJourneyTaskV1", func() {
			Context("Create journey", func() {
				It("should return empty result with calling producer", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.CreateJourney,
						Value:       journeysTable[0],
					}).Times(1)
					mockMetrics.EXPECT().CreateJourneyCounterInc().Times(1)

					result, err := api.CreateJourneyTaskV1(ctx, &desc.CreateJourneyTaskRequestV1{
						UserId:    journeysTable[0].UserID,
						Address:   journeysTable[0].Address,
						StartTime: timestamppb.New(journeysTable[0].StartTime),
						EndTime:   timestamppb.New(journeysTable[0].EndTime),
					})

					Expect(result).Should(Equal(&emptypb.Empty{}))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journey in request", func() {
				It("should return error without calling producer", func() {
					mockProducer.EXPECT().Send(gomock.Any()).Times(0)

					result, err := api.CreateJourneyTaskV1(ctx, &desc.CreateJourneyTaskRequestV1{
						UserId:    0,
						Address:   journeysTable[0].Address,
						StartTime: timestamppb.New(journeysTable[0].StartTime),
						EndTime:   timestamppb.New(journeysTable[0].EndTime),
					})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in producer", func() {
				It("should return error", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.CreateJourney,
						Value:       journeysTable[0],
					}).Times(1).Return(errProducer)

					result, err := api.CreateJourneyTaskV1(ctx, &desc.CreateJourneyTaskRequestV1{
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

		Context("MultiCreateJourneyTaskV1", func() {
			var journeysTableZeroIds []models.Journey

			BeforeEach(func() {
				journeysTableZeroIds = make([]models.Journey, len(journeysTable))
				copy(journeysTableZeroIds, journeysTable)
				for i := range journeysTableZeroIds {
					journeysTableZeroIds[i].JourneyID = 0
				}
			})

			Context("Success create journeys", func() {
				It("should return success empty result with calling producer", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.MultiCreateJourney,
						Value:       journeysTableZeroIds[0:2],
					}).Times(1)
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.MultiCreateJourney,
						Value:       journeysTableZeroIds[2:],
					}).Times(1)
					mockMetrics.EXPECT().MultiCreateJourneyCounterInc().Times(1)

					req := &desc.MultiCreateJourneyTaskRequestV1{}
					for _, journey := range journeysTable {
						req.Journeys = append(req.Journeys, &desc.CreateJourneyRequestV1{
							UserId:    journey.UserID,
							Address:   journey.Address,
							StartTime: timestamppb.New(journey.StartTime),
							EndTime:   timestamppb.New(journey.EndTime),
						})
					}

					result, err := api.MultiCreateJourneyTaskV1(ctx, req)

					Expect(result).Should(Equal(&emptypb.Empty{}))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journeys count in request", func() {
				It("should return error without calling producer", func() {
					mockProducer.EXPECT().Send(gomock.Any()).Times(0)

					result, err := api.MultiCreateJourneyTaskV1(ctx, &desc.MultiCreateJourneyTaskRequestV1{})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in producer", func() {
				It("should return error", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.MultiCreateJourney,
						Value:       journeysTableZeroIds[0:2],
					}).Times(1)
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.MultiCreateJourney,
						Value:       journeysTableZeroIds[2:],
					}).Return(errProducer).Times(1)

					req := &desc.MultiCreateJourneyTaskRequestV1{}
					for _, journey := range journeysTable {
						req.Journeys = append(req.Journeys, &desc.CreateJourneyRequestV1{
							UserId:    journey.UserID,
							Address:   journey.Address,
							StartTime: timestamppb.New(journey.StartTime),
							EndTime:   timestamppb.New(journey.EndTime),
						})
					}

					result, err := api.MultiCreateJourneyTaskV1(ctx, req)

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Context("RemoveJourneyTaskV1", func() {
			Context("Success remove journey with calling producer", func() {
				It("should return empty result", func() {
					journeyID := uint64(1)
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.DeleteJourney,
						Value:       journeyID,
					}).Times(1)
					mockMetrics.EXPECT().DeleteJourneyCounterInc().Times(1)

					result, err := api.RemoveJourneyTaskV1(ctx, &desc.RemoveJourneyTaskRequestV1{JourneyId: journeyID})

					Expect(result).Should(Equal(&emptypb.Empty{}))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journeyId in request", func() {
				It("should return error without calling producer", func() {
					journeyID := uint64(0)
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.DeleteJourney,
						Value:       journeyID,
					}).Times(0)

					result, err := api.RemoveJourneyTaskV1(ctx, &desc.RemoveJourneyTaskRequestV1{JourneyId: journeyID})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in producer", func() {
				It("should return error", func() {
					mockProducer.EXPECT().Send(gomock.Any()).Return(errProducer).Times(1)

					result, err := api.RemoveJourneyTaskV1(ctx, &desc.RemoveJourneyTaskRequestV1{JourneyId: 1})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})
		})

		Context("UpdateJourneyTaskV1", func() {
			Context("Success update journey", func() {
				It("should return success empty result", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.UpdateJourney,
						Value:       journeysTable[2],
					}).Times(1)
					mockMetrics.EXPECT().UpdateJourneyCounterInc().Times(1)

					result, err := api.UpdateJourneyTaskV1(ctx, &desc.UpdateJourneyTaskRequestV1{
						Journey: &desc.Journey{
							JourneyId: journeysTable[2].JourneyID,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(Equal(&emptypb.Empty{}))
					Expect(err).Should(BeNil())
				})
			})

			Context("Incorrect journey id in request", func() {
				It("should return error without calling producer", func() {
					mockProducer.EXPECT().Send(gomock.Any()).Times(0)

					result, err := api.UpdateJourneyTaskV1(ctx, &desc.UpdateJourneyTaskRequestV1{
						Journey: &desc.Journey{
							JourneyId: 0,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})

			Context("Error in producer", func() {
				It("should return error", func() {
					mockProducer.EXPECT().Send(kafka.Message{
						MessageType: kafka.UpdateJourney,
						Value:       journeysTable[2],
					}).Return(errProducer).Times(1)

					result, err := api.UpdateJourneyTaskV1(ctx, &desc.UpdateJourneyTaskRequestV1{
						Journey: &desc.Journey{
							JourneyId: journeysTable[2].JourneyID,
							UserId:    journeysTable[2].UserID,
							Address:   journeysTable[2].Address,
							StartTime: timestamppb.New(journeysTable[2].StartTime),
							EndTime:   timestamppb.New(journeysTable[2].EndTime),
						},
					})

					Expect(result).Should(BeNil())
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

})
