package kafka

// MessageType - represents type for iota with Create, MultiCreate, Update, Delete operations
type MessageType int

const (
	// CreateJourney - Create Journey via Producer
	CreateJourney MessageType = iota
	// MultiCreateJourney - Create several Journeys via Producer
	MultiCreateJourney
	// UpdateJourney - Update Journey via Producer
	UpdateJourney
	// DeleteJourney - Delete Journey via Producer
	DeleteJourney
)

// Message - message for Kafka
type Message struct {
	MessageType MessageType
	Value       interface{}
}
