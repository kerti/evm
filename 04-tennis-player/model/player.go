package model

import (
	"math/rand"

	"github.com/gofrs/uuid"
	"github.com/kerti/evm/04-tennis-player/util/failure"
)

// Player represents a Player entity object
type Player struct {
	ID          uuid.UUID   `json:"id" db:"entity_id" validate:"min=36,max=36"`
	Name        string      `json:"name" db:"name"`
	ReadyToPlay bool        `json:"readyToPlay" db:"ready_to_play"`
	Containers  []Container `json:"containers" db:"-"`
}

// NewPlayerFromInput creates a new Player from its input object
func NewPlayerFromInput(input PlayerInput) Player {
	id := input.ID
	if input.ID == uuid.Nil {
		id, _ = uuid.NewV4()
	}
	return Player{
		ID:          id,
		Name:        input.Name,
		ReadyToPlay: false,
	}
}

// AttachContainers attaches containers to a player
func (p *Player) AttachContainers(containers []Container) Player {
	for _, container := range containers {
		if container.PlayerID == p.ID {
			p.Containers = append(p.Containers, container)
		}
	}
	return *p
}

// AddBall adds a single ball randomly into one of the player's containers
func (p *Player) AddBall() error {
	if err := p.ValidateAddBall(); err != nil {
		return err
	}

	randomContainerIndex := rand.Intn(len(p.Containers))

	containers := make([]Container, 0)
	for idx, container := range p.Containers {
		if idx == randomContainerIndex {
			container.AddBall()

			if container.IsFull() {
				p.ReadyToPlay = true
			}
		}

		if container.IsFull() {
			p.ReadyToPlay = true
		}

		containers = append(containers, container)
	}
	p.Containers = containers

	return nil
}

// ValidateAddBall checks if the player can still add balls to one of his containers
func (p *Player) ValidateAddBall() error {
	if p.ReadyToPlay {
		return failure.OperationNotPermitted("addBall", "player", "the player is ready to play")
	}

	if len(p.Containers) == 0 {
		return failure.OperationNotPermitted("addBall", "player", "the player has no containers to put the ball into")
	}

	for _, container := range p.Containers {
		if container.IsFull() {
			return failure.OperationNotPermitted("addBall", "player", "the player should already be ready to play")
		}
	}

	return nil
}

// PlayerInput represents the input object for creating new Players
type PlayerInput struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name"`
}

// PlayerAddBallInput represents the input object for players to add balls
type PlayerAddBallInput struct {
	PlayerID uuid.UUID `json:"playerId"`
}
