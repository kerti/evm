package service

import (
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/repository"
	"github.com/kerti/evm/04-tennis-player/util/failure"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

// Container is the service provider interface
type Container interface {
	Startup()
	Shutdown()
	Create(input model.ContainerInput) (*model.Container, error)
}

// ContainerImpl is the service provider implementation
type ContainerImpl struct {
	DB                  *database.MySQL      `inject:"mysql"`
	ContainerRepository repository.Container `inject:"containerRepository"`
	PlayerRepository    repository.Player    `inject:"playerRepository"`
}

// Startup performs startup functions
func (s *ContainerImpl) Startup() {
	logger.Trace("Container Service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *ContainerImpl) Shutdown() {
	logger.Trace("Container Service shutting down...")
}

// Create creates a new Container
func (s *ContainerImpl) Create(input model.ContainerInput) (*model.Container, error) {
	container := model.NewContainerFromInput(input)

	playerExists, err := s.PlayerRepository.ExistsByID(container.PlayerID)
	if err != nil {
		return nil, err
	}

	if !playerExists {
		return nil, failure.OperationNotPermitted("create", "Container", "specified Player does not exist")
	}

	err = s.ContainerRepository.Create(container)
	return &container, err
}
