package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"time"
)

type CommentDto struct {
	Id        uuid.UUID `json:"id"        bson:"id"`
	Text      string    `json:"text"        bson:"text"`
	Locale    string    `json:"locale"        bson:"locale"`
	CreatedAt time.Time `json:"createdAt"        bson:"createdAt"`
	CreatedBy string    `json:"createdBy"        bson:"createdDate"`
}

type CommentedEntity struct {
	Id         uuid.UUID     `json:"id"        bson:"id"`
	EntityType string        `json:"entityType"        bson:"entityType"`
	EntityId   string        `json:"entityId"        bson:"entityId"`
	Comments   []*CommentDto `json:"comments"        bson:"comments"`
	Version    int           `json:"version"    bson:"version"`
	UpdatedAt  time.Time     `json:"updatedAt"        bson:"updatedAt"`
}

var _ = eh.Entity(&CommentedEntity{})
var _ = eh.Versionable(&CommentedEntity{})

// EntityID implements the EntityID method of the eventhorizon.Entity interface.
func (t *CommentedEntity) EntityID() uuid.UUID {
	return t.Id
}

// AggregateVersion implements the AggregateVersion method of the
// eventhorizon.Versionable interface.
func (t *CommentedEntity) AggregateVersion() int {
	return t.Version
}
