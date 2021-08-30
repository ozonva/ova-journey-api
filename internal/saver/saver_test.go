//go:build !race
// +build !race

package saver

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozonva/ova-journey-api/internal/mocks"
	"github.com/ozonva/ova-journey-api/internal/models"
	"sync"
	"time"
)

var _ = Describe("Saver", func() {
	var (
		ctrl          *gomock.Controller
		mockFlusher   *mocks.MockFlusher
		s             Saver
		capacity      uint
		delayFlushing time.Duration
	)

	journeysTable := []models.Journey{
		{JourneyID: 0, UserID: 1, Address: "Воронеж", Description: ""},
		{JourneyID: 1, UserID: 1, Address: "Уфа", Description: ""},
		{JourneyID: 2, UserID: 2, Address: "Москва", Description: ""},
		{JourneyID: 3, UserID: 2, Address: "Лондон", Description: ""},
		{JourneyID: 4, UserID: 3, Address: "Новосибирск", Description: ""},
	}

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("building saver with incorrect parameters", func() {
		When("capacity is incorrect", func() {
			BeforeEach(func() {
				delayFlushing = time.Second
			})
			It("should panic", func() {
				Expect(func() { NewSaver(0, mockFlusher, delayFlushing) }).Should(Panic())
			})
		})
		When("flusher is incorrect", func() {
			BeforeEach(func() {
				capacity = 2
				delayFlushing = time.Second
			})
			It("should panic", func() {
				Expect(func() { NewSaver(capacity, nil, delayFlushing) }).Should(Panic())
			})
		})
		When("flushing delay is incorrect", func() {
			BeforeEach(func() {
				capacity = 2
			})
			It("should panic", func() {
				Expect(func() { NewSaver(capacity, mockFlusher, -1) }).Should(Panic())
				Expect(func() { NewSaver(capacity, mockFlusher, 0) }).Should(Panic())
			})
		})
	})

	Context("using correct built saver", func() {
		JustBeforeEach(func() {
			s = NewSaver(capacity, mockFlusher, delayFlushing)
		})

		Context("using closed saver", func() {
			BeforeEach(func() {
				capacity = 2
				delayFlushing = time.Hour
			})

			When("try to close on already closed saver", func() {
				It("should return ErrSaverIsClosed without attempts to flush", func() {
					mockFlusher.EXPECT().Flush(gomock.Any()).Times(0)

					closeResult := s.Close()
					Expect(closeResult).Should(BeNil())

					closeResult = s.Close()
					Expect(closeResult).Should(Equal(ErrSaverIsClosed))
				})
			})

			When("try to save on closed saver", func() {
				It("should return ErrSaverIsClosed without attempts to flush", func() {
					mockFlusher.EXPECT().Flush(gomock.Any()).Times(0)

					closeResult := s.Close()
					Expect(closeResult).Should(BeNil())

					for _, journey := range journeysTable {
						result := s.Save(journey)
						Expect(result).Should(Equal(ErrSaverIsClosed))
					}
				})
			})

			When("try to save on closed saver after add one journey with success flushing", func() {
				It("should return ErrSaverIsClosed and try to flush 1 item", func() {
					mockFlusher.EXPECT().Flush(journeysTable[:1]).Times(1).Return(nil)

					saveResult := s.Save(journeysTable[0])
					closeResult := s.Close()
					Expect(saveResult).Should(BeNil())
					Expect(closeResult).Should(BeNil())

					for _, journey := range journeysTable {
						result := s.Save(journey)
						Expect(result).Should(Equal(ErrSaverIsClosed))
					}
				})
			})

			When("try to close saver after add one journey with failed flushing", func() {
				It("should return ErrPartOfDataIsNotFlushed for close", func() {
					mockFlusher.EXPECT().Flush(journeysTable[:1]).Times(1).Return(journeysTable[:1])

					saveResult := s.Save(journeysTable[0])
					closeResult := s.Close()

					Expect(saveResult).Should(BeNil())
					Expect(closeResult).Should(Equal(ErrPartOfDataIsNotFlushed))
				})
			})
		})

		Context("periodic flushing", func() {
			BeforeEach(func() {
				capacity = 2
				delayFlushing = time.Millisecond * 500
			})
			When("delay between Save() calls is greater then delay between flushing", func() {
				It("should call flusher several times (all data size/ saver capacity) without errors", func() {
					flusherCallsLimit := len(journeysTable) / int(capacity)

					if len(journeysTable)%int(capacity) > 0 {
						flusherCallsLimit++
					}
					fmt.Println(flusherCallsLimit)
					mockFlusher.EXPECT().Flush(gomock.Any()).MinTimes(flusherCallsLimit)

					wg := sync.WaitGroup{}
					wg.Add(1)

					go func() {
						dataSendTicker := time.NewTicker(delayFlushing + delayFlushing/2)
						defer dataSendTicker.Stop()
						i := 0
						for {
							<-dataSendTicker.C
							result := s.Save(journeysTable[i])
							Expect(result).Should(BeNil())
							i++
							if i >= len(journeysTable) {
								wg.Done()
								return
							}
						}
					}()

					wg.Wait()
				})
			})
			When("delay between Save() calls is less then delay between flushing", func() {
				It("should call flusher several times and return at least one ErrInternalBufferIsFull", func() {
					flusherCallsLimit := len(journeysTable) / int(capacity)

					if len(journeysTable)%int(capacity) > 0 {
						flusherCallsLimit++
					}
					mockFlusher.EXPECT().Flush(gomock.Any()).MaxTimes(flusherCallsLimit)

					results := make([]error, 0, len(journeysTable))

					wg := sync.WaitGroup{}
					wg.Add(1)

					go func() {
						dataSendTicker := time.NewTicker(delayFlushing / time.Duration(capacity*2))
						defer dataSendTicker.Stop()
						i := 0
						for {
							<-dataSendTicker.C
							results = append(results, s.Save(journeysTable[i]))
							i++
							if i >= len(journeysTable) {
								wg.Done()
								return
							}
						}
					}()

					wg.Wait()

					Expect(results).Should(ContainElements(ErrInternalBufferIsFull))
				})
			})
		})
	})

})
