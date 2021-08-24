package saver

import (
	"errors"
	"github.com/ozonva/ova-journey-api/internal/flusher"
	"github.com/ozonva/ova-journey-api/internal/models"
	"sync"
	"time"
)

// InternalBufferIsFullError - this error occurs when you try to add a new element using the Saver.Save() method
// if there are already capacity elements in the internal buffer
var InternalBufferIsFullError = errors.New("cannot add new journey, saver internal buffer is full")

// IsClosedError - this error occurs when you try to add a new element using the Saver.Save() method after calling
// Saver.Close() method. closed Saver cannot be used, you should create a new Saver
var IsClosedError = errors.New("saver is closed and cannot be used anymore")

// PartOfDataIsNotFlushedError - this error occurs when part of the data remains unwritten
// from internal buffer to the storage by flusher
var PartOfDataIsNotFlushedError = errors.New("part of journeys was not flushed")

// Saver represents the object used for saving journeys in storage
type Saver interface {
	// Save - add new models.Journey for saving
	Save(entity models.Journey) error
	// Close - close the Saver
	Close() error
}

type saverState int

const (
	closed saverState = iota
	run
)

type saver struct {
	sync.Mutex
	flusher flusher.Flusher
	buffer  []models.Journey
	done    chan struct{}
	state   saverState
}

// Save - add new journey to internal buffer of Saver
func (s *saver) Save(journey models.Journey) error {
	// lock the state and internal buffer
	s.Lock()
	defer s.Unlock()

	if s.state == closed {
		return IsClosedError
	}

	if len(s.buffer) == cap(s.buffer) {
		return InternalBufferIsFullError
	}

	s.buffer = append(s.buffer, journey)
	return nil
}

// Close - close the Saver with flushing all remain data from internal buffer.
func (s *saver) Close() error {
	// lock the state
	s.Lock()
	if s.state == closed {
		return IsClosedError
	}

	s.state = closed
	s.done <- struct{}{}
	close(s.done)
	s.Unlock()

	return s.uploadToFlusher()
}

func (s *saver) uploadToFlusher() error {
	// lock internal buffer
	s.Lock()
	defer s.Unlock()

	if len(s.buffer) == 0 {
		return nil
	}

	saveFailedJourneys := s.flusher.Flush(s.buffer)
	s.buffer = s.buffer[:0]

	// if not all data was flushed, restore them to try flush again in next call
	if len(saveFailedJourneys) > 0 {
		s.buffer = append(s.buffer, saveFailedJourneys...)
		return PartOfDataIsNotFlushedError
	}
	return nil
}

// NewSaver return Saver with periodic flushing Journeys data to the storage using flusher.
//
// Use Saver.Save() method to add new journey for flushing.
// For collecting data between flushing attempts used internal buffer with capacity size.
// If internal buffer is full the Saver.Save() method returns InternalBufferIsFullError without adding journey for flushing.
// If Saver is already closed the Saver.Save() method returns IsClosedError without trying to flush journey.
//
// Use Saver.Close() method to immediately flush data and close Saver.
// After closing Saver stops to try flushing data and cannot be used anymore.
// If not all data from internal buffer was flushed the Saver.Close() method returns PartOfDataIsNotFlushedError.
func NewSaver(
	capacity uint,
	flusher flusher.Flusher,
	delayBetweenFlushing time.Duration,
) Saver {
	if capacity < 1 {
		panic("capacity must be greater then 0")
	}
	if flusher == nil {
		panic("flusher cannot be nil")
	}
	if delayBetweenFlushing < 1 {
		panic("delayBetweenFlushing must be greater then 0")
	}

	s := saver{
		flusher: flusher,
		done:    make(chan struct{}),
		buffer:  make([]models.Journey, 0, capacity),
		state:   run,
	}

	go func() {
		ticker := time.NewTicker(delayBetweenFlushing)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = s.uploadToFlusher()
			case <-s.done:
				return
			}
		}
	}()

	return &s
}
