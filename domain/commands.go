package domain

import (
	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
	_ "google.golang.org/api/monitoring/v1"
	"time"
	_ "time"
)

func init() {
	eh.RegisterCommand(func() eh.Command { return &CreateCommentCommand{} })
}

const (
	// CreateCommand is the type for the Create command.
	CreateCommand = eh.CommandType("comment:create")
)

// Static type check that the eventhorizon.Command interface is implemented.
var _ = eh.Command(&CreateCommentCommand{})

// Create creates a new todo list.
type CreateCommentCommand struct {
	ID               uuid.UUID
	InteractableType string
	InteractableId   string
	Text             string
	Locale           string
	CreatedAt        time.Time
	CreatedBy        string
}

func (c *CreateCommentCommand) AggregateType() eh.AggregateType { return AggregateType }
func (c *CreateCommentCommand) AggregateID() uuid.UUID          { return c.ID }
func (c *CreateCommentCommand) CommandType() eh.CommandType     { return CreateCommand }
