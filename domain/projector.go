package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
)

// Projector is a projector of comment events on the CommentedEntity read model.
type Projector struct{

	QueryRepo CommentedEntityRepository
}

// ProjectorType implements the ProjectorType method of the
// eventhorizon.Projector interface.
func (p *Projector) ProjectorType() projector.Type {
	return projector.Type(string(AggregateType) + "_projector")
}

// Project implements the Project method of the eventhorizon.Projector interface.
func (p *Projector) Project(ctx context.Context,
	event eh.Event, entity eh.Entity) (eh.Entity, error) {

	model, ok := entity.(*CommentedEntity)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

	switch event.EventType() {
	case CommentCreatedEvent:

		// get entity if exists if not create a new one
		commentedEntity, err = p.QueryRepo.FindByEntityTypeAndId(ctx,model.EntityType)

		data, ok := event.Data().(*CommentCreated)
		if !ok {
			return nil, errors.New("invalid event data")
		}
		model.Id, _ = uuid.Parse("26fa44f2-5c53-4616-be1d-571e818264f1")
		model.EntityId = data.EntityId
		model.EntityType = data.EntityType
		model.Comments = append(model.Comments, &CommentDto{
			Id:        data.Id,
			Text:      data.Text,
			Locale:    data.Locale,
			CreatedAt: data.CreatedAt,
			CreatedBy: data.CreatedBy,
		})

	default:
		// Also return the modele here to not delete it.
		return model, fmt.Errorf("could not project event: %s", event.EventType())
	}

	// Always increment the version and set update time on successful updates.
	model.Version++
	model.UpdatedAt = TimeNow()
	return model, nil
}
