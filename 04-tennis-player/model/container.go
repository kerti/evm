package model

import (
	"github.com/gofrs/uuid"
	"github.com/kerti/evm/04-tennis-player/util/failure"
)

// Container represents a Container entity object
type Container struct {
	ID        uuid.UUID `json:"id" db:"entity_id" validate:"min=36,max=36"`
	PlayerID  uuid.UUID `json:"playerId" db:"player_entity_id" validate:"min=36,max=36"`
	Capacity  int       `json:"capacity" db:"capacity" validate:"min=0"`
	BallCount int       `json:"ballCount" db:"ball_count" validate:"min=0"`
}

// NewContainerFromInput creates a new Container from its input object
func NewContainerFromInput(input ContainerInput) Container {
	id := input.ID
	if input.ID == uuid.Nil {
		id, _ = uuid.NewV4()
	}
	return Container{
		ID:        id,
		PlayerID:  input.PlayerID,
		Capacity:  input.Capacity,
		BallCount: 0,
	}
}

// AddBall adds a single ball into a Container
func (c *Container) AddBall() (Container, error) {
	if c.IsFull() {
		return *c, failure.OperationNotPermitted("addBall", "container", "the container is full")
	}

	c.BallCount++

	return *c, nil
}

// IsFull checks whether a ball can be added into a Container
func (c *Container) IsFull() bool {
	return c.Capacity == c.BallCount
}

// ContainerInput represents the input object for creating new Containers
type ContainerInput struct {
	ID       uuid.UUID `json:"id,omitempty"`
	PlayerID uuid.UUID `json:"playerId"`
	Capacity int       `json:"capacity"`
}
