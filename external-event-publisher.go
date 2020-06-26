package main

import (
	"context"
	eh "github.com/looplab/eventhorizon"
	"log"
)

// public domain events externally

// Logger is a simple event handler for logging all events.
type ExternalEventPublisher struct{}

// HandlerType implements the HandlerType method of the eventhorizon.EventHandler interface.
func (e *ExternalEventPublisher) HandlerType() eh.EventHandlerType {
	return "externalEventPublisher"
}

// HandleEvent implements the HandleEvent method of the EventHandler interface.
func (e *ExternalEventPublisher) HandleEvent(ctx context.Context, event eh.Event) error {

	// publish the message to pub sub


	log.Printf("EVENT %s published", event)
	return nil
}