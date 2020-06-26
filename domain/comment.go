package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"time"
)

func init() {
	eh.RegisterAggregate(func(id uuid.UUID) eh.Aggregate {
		return &Comment{
			AggregateBase: events.NewAggregateBase(AggregateType, id),
		}
	})
}

// AggregateType is the aggregate type for the interactable
const AggregateType = eh.AggregateType("comment")

// Aggregate is an aggregate for a comment
type Comment struct {
	*events.AggregateBase

	id         uuid.UUID
	entityType string
	entityId   string
	text       string
	locale     string
	createdAt  time.Time
	createdBy  string
}

// TimeNow is a mockable version of time.Now.
var TimeNow = time.Now

// HandleCommand implements the HandleCommand method of the
// eventhorizon.CommandHandler interface.
func (a *Comment) HandleCommand(ctx context.Context, cmd eh.Command) error {
	switch cmd := cmd.(type) {
	case *CreateCommentCommand:
		a.StoreEvent(CommentCreatedEvent, &CommentCreated{
			Id:          cmd.AggregateID(),
			EntityType:  cmd.InteractableType,
			EntityId:    cmd.InteractableId,
			Text:        cmd.Text,
			Locale:      cmd.Locale,
			CreatedAt: cmd.CreatedAt,
			CreatedBy:   cmd.CreatedBy,
		}, time.Now())
	default:
		return fmt.Errorf("could not handle command: %s", cmd.CommandType())
	}

	return nil
}

// ApplyEvent implements the ApplyEvent method of the
// eventhorizon.Aggregate interface.
func (a *Comment) ApplyEvent(ctx context.Context, event eh.Event) error {
	switch event.EventType() {
	case CommentCreatedEvent:
		data, ok := event.Data().(*CommentCreated)
		if !ok {
			return errors.New("invalid event data")
		}
		a.id = data.Id
		a.entityId = data.EntityId
		a.entityType = data.EntityType
		a.text = data.Text
		a.locale = data.Locale
		a.createdAt = data.CreatedAt
		a.createdBy = data.CreatedBy
	}

	return nil
}
