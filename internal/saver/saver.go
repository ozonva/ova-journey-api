package saver

import (
	"errors"
	"github.com/ozonva/ova-journey-api/internal/flusher"
	"github.com/ozonva/ova-journey-api/internal/models"
	"sync"
	"time"
)

var InternalBufferIsFullError = errors.New("cannot add new journey, saver internal buffer is full")
var IsClosedError = errors.New("saver is closed and cannot be used anymore")
var IsClosedWithRemainDataError = errors.New("saver is closed, but part of journeys was not flushed")

type Saver interface {
	Save(entity models.Journey) error
	Close() error
}

type saverState int

const (
	Closed saverState = iota
	Run
)

type saver struct {
	sync.Mutex
	flusher flusher.Flusher
	buffer  []models.Journey
	close   chan struct{}
	done    chan struct{}
	state   saverState
}

func (s *saver) Save(journey models.Journey) error {
	s.Lock()
	defer s.Unlock()

	if s.state == Closed {
		return IsClosedError
	}

	if len(s.buffer) == cap(s.buffer) {
		return InternalBufferIsFullError
	}

	s.buffer = append(s.buffer, journey)
	return nil
}

func (s *saver) Close() error {
	s.close <- struct{}{}
	<-s.done
	close(s.close)

	if len(s.buffer) > 0 {
		return IsClosedWithRemainDataError
	}
	return nil
}

func (s *saver) uploadToFlusher() {
	saveFailedJourneys := s.flusher.Flush(s.buffer)
	s.buffer = s.buffer[:0]
	// if not all data was flushed, restore them to try flush again in next call
	if len(saveFailedJourneys) > 0 {
		s.buffer = append(s.buffer, saveFailedJourneys...)
	}
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	delayBetweenFlushing time.Duration,
) Saver {
	s := saver{
		flusher: flusher,
		close:   make(chan struct{}),
		done:    make(chan struct{}),
		buffer:  make([]models.Journey, 0, capacity),
		state:   Run,
	}

	go func() {
		ticker := time.NewTicker(delayBetweenFlushing)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.Lock()
				s.uploadToFlusher()
				s.Unlock()
			case <-s.close:
				s.Lock()
				if len(s.buffer) > 0 {
					s.uploadToFlusher()
				}
				s.state = Closed
				s.Unlock()

				s.done <- struct{}{}
				close(s.done)
				return
			}
		}
	}()

	return &s
}
