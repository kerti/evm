package service

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/repository"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

// Player is the service provider interface
type Player interface {
	Startup()
	Shutdown()
	AddBall(playerID uuid.UUID) (*model.Player, error)
}

// PlayerImpl is the service provider implementation
type PlayerImpl struct {
	DB                  *database.MySQL      `inject:"mysql"`
	ContainerRepository repository.Container `inject:"containerRepository"`
	PlayerRepository    repository.Player    `inject:"playerRepository"`
	mux                 sync.Mutex
}

// Startup performs startup functions
func (s *PlayerImpl) Startup() {
	logger.Trace("Player Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *PlayerImpl) Shutdown() {
	logger.Trace("Player Service shutting down...")
}

// AddBall adds a ball into one of the player's containers
func (s *PlayerImpl) AddBall(playerID uuid.UUID) (*model.Player, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	player, err := s.PlayerRepository.ResolveByID(playerID)
	if err != nil {
		return nil, err
	}

	containers, err := s.ContainerRepository.ResolveByPlayerID(player.ID)
	if err != nil {
		return nil, err
	}

	player.AttachContainers(containers)

	err = player.AddBall()
	if err != nil {
		return nil, err
	}

	err = s.DB.WithTransaction(s.DB, func(tx *sqlx.Tx, e chan error) {
		if err := s.PlayerRepository.TxUpdate(tx, *player); err != nil {
			e <- err
			return
		}

		if err := s.ContainerRepository.TxBulkUpdate(tx, player.Containers); err != nil {
			e <- err
			return
		}

		e <- nil
	})

	return player, err
}
