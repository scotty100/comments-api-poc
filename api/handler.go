package api

import (
	"comments-api/domain"
	"context"
	"fmt"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/aggregatestore/events"
	"github.com/looplab/eventhorizon/commandhandler/aggregate"
	eventbus "github.com/looplab/eventhorizon/eventbus/local"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	eventstore "github.com/looplab/eventhorizon/eventstore/mongodb"
	"github.com/looplab/eventhorizon/httputils"
	repo "github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/looplab/eventhorizon/repo/version"
	"log"
	"net/http"
	"os"
)

// Handler is a http.Handler for the TodoMVC app.
type Handler struct {
	http.Handler

	EventBus       eh.EventBus
	CommandHandler eh.CommandHandler
	Repo           eh.ReadWriteRepo
}

// Logger is a simple event handler for logging all events.
type Logger struct{}

// HandlerType implements the HandlerType method of the eventhorizon.EventHandler interface.
func (l *Logger) HandlerType() eh.EventHandlerType {
	return "logger"
}

// HandleEvent implements the HandleEvent method of the EventHandler interface.
func (l *Logger) HandleEvent(ctx context.Context, event eh.Event) error {
	log.Printf("EVENT %s", event)
	return nil
}

// NewHandler sets up the full Event Horizon domain for the Comments app and
// returns a handler exposing some of the components.
func NewHandler() (*Handler, error) {
	// Use MongoDB in Docker with fallback to localhost.
	dbURL := os.Getenv("MONGO_HOST")
	if dbURL == "" {
		dbURL = "localhost:27017"
	}
	dbURL = "mongodb://" + dbURL

	// Create the event store.
	eventStore, err := eventstore.NewEventStore(dbURL, "comments")
	if err != nil {
		return nil, fmt.Errorf("could not create event store: %s", err)
	}

	// Create the event bus that distributes events.
	eventBus := eventbus.NewEventBus(nil)
	go func() {
		for e := range eventBus.Errors() {
			log.Printf("eventbus: %s", e.Error())
		}
	}()

	// Add a logger as an observer.
	eventBus.AddObserver(eh.MatchAny(), &Logger{})

	// Create the aggregate repository.
	aggregateStore, err := events.NewAggregateStore(eventStore, eventBus)
	if err != nil {
		return nil, fmt.Errorf("could not create aggregate store: %s", err)
	}

	// Create the aggregate command handler.
	aggregateCommandHandler, err := aggregate.NewCommandHandler(domain.AggregateType, aggregateStore)
	if err != nil {
		return nil, fmt.Errorf("could not create command handler: %s", err)
	}

	// Create a tiny logging middleware for the command handler.
	commandHandlerLogger := func(h eh.CommandHandler) eh.CommandHandler {
		return eh.CommandHandlerFunc(func(ctx context.Context, cmd eh.Command) error {
			log.Printf("CMD %#v", cmd)
			return h.HandleCommand(ctx, cmd)
		})
	}
	commandHandler := eh.UseCommandHandlerMiddleware(aggregateCommandHandler, commandHandlerLogger)

	// Create the repository and wrap in a version repository.
	repo, err := repo.NewRepo(dbURL, "comments", "commentedEntities")
	if err != nil {
		return nil, fmt.Errorf("could not create invitation repository: %s", err)
	}
	repo.SetEntityFactory(func() eh.Entity { return &domain.CommentedEntity{} })
	commentedRepo := version.NewRepo(repo)

	// Create the read model projector.
	projector := projector.NewEventHandler(&domain.Projector{}, commentedRepo)
	projector.SetEntityFactory(func() eh.Entity { return &domain.CommentedEntity{} })
	eventBus.AddHandler(eh.MatchAnyEventOf(
		domain.CommentCreatedEvent,
	), projector)

	// Handle the API.
	h := http.NewServeMux()
	h.Handle("/events/", httputils.EventBusHandler(eventBus, eh.MatchAny(), "any"))
	//h.Handle("/api/", httputils.QueryHandler(commentedRepo))
	h.Handle("/comments", httputils.CommandHandler(commandHandler, domain.CreateCommand))

	// Simple HTTP request logging middleware as final handler.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})

	return &Handler{
		Handler:        handler,
		EventBus:       eventBus,
		CommandHandler: commandHandler,
		Repo:           commentedRepo,
	}, nil

}
