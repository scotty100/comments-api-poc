package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"time"
)

const (
	// CommentCreatedEvent is the event after a comment is created.
	CommentCreatedEvent = eh.EventType("comment:created")
)

func init() {
	eh.RegisterEventData(CommentCreatedEvent, func() eh.EventData {
		return &CommentCreated{}
	})
}

type CommentCreated struct {
	Id         uuid.UUID `json:"id"     bson:"id"`
	EntityType string    `json:"entityType"     bson:"entityType"`
	EntityId   string    `json:"entityId"     bson:"entityId"`
	Text       string    `json:"text"     bson:"text"`
	Locale     string    `json:"locale"     bson:"locale"`
	CreatedAt  time.Time `json:"createdAt"     bson:"createdAt"`
	CreatedBy  string    `json:"createdBy"     bson:"createdBy"`
}
