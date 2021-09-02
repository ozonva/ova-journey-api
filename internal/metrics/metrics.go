package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics - interface for Create/Update/Delete operations on journey
type Metrics interface {
	CreateJourneyCounterInc()
	MultiCreateJourneyCounterInc()
	UpdateJourneyCounterInc()
	DeleteJourneyCounterInc()
}

type metrics struct {
	createJourneySuccessCounter      prometheus.Counter
	multiCreateJourneySuccessCounter prometheus.Counter
	updateJourneySuccessCounter      prometheus.Counter
	deleteJourneySuccessCounter      prometheus.Counter
}

// NewMetrics - creates new Metrics object
func NewMetrics(namespace, subsystem string) Metrics {
	return &metrics{
		createJourneySuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "create_count_total",
			Help:      "Total count of successful requests to create journey",
		}),
		multiCreateJourneySuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "multi_create_count_total",
			Help:      "Total count of successful requests to chunked create journeys",
		}),
		updateJourneySuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "update_count_total",
			Help:      "Total count of successful requests to update journey",
		}),
		deleteJourneySuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "delete_count_total",
			Help:      "Total count of successful requests to delete journey",
		}),
	}
}

func (m *metrics) CreateJourneyCounterInc() {
	m.createJourneySuccessCounter.Inc()
}

func (m *metrics) MultiCreateJourneyCounterInc() {
	m.multiCreateJourneySuccessCounter.Inc()
}

func (m *metrics) UpdateJourneyCounterInc() {
	m.updateJourneySuccessCounter.Inc()
}

func (m *metrics) DeleteJourneyCounterInc() {
	m.deleteJourneySuccessCounter.Inc()
}
