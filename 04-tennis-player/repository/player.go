package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/util/failure"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

const (
	queryInsertPlayer = `
		INSERT INTO players (
			players.entity_id,
			players.name,
			players.ready_to_play
		) VALUES (
			:entity_id,
			:name,
			:ready_to_play)`

	querySelectPlayer = `
		SELECT
			players.entity_id,
			players.name,
			players.ready_to_play
		FROM players`

	queryUpdatePlayer = `
		UPDATE players
		SET
			name = :name,
			ready_to_play = :ready_to_play
		WHERE entity_id = :entity_id`
)

// Player is the Player repository interface
type Player interface {
	Startup()
	Shutdown()
	ExistsByID(id uuid.UUID) (exists bool, err error)
	Create(player model.Player) (err error)
	ResolveByID(id uuid.UUID) (player *model.Player, err error)
	TxUpdate(tx *sqlx.Tx, player model.Player) (err error)
}

// PlayerMySQLRepo is the repository for Players implemented with MySQL backend
type PlayerMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup performs startup functions
func (r *PlayerMySQLRepo) Startup() {
	logger.Trace("Player Repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *PlayerMySQLRepo) Shutdown() {
	logger.Trace("Player Repository starting up...")
}

// ExistsByID checks whether a Player exists by its ID
func (r *PlayerMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM players WHERE players.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// Create creates a new Player
func (r *PlayerMySQLRepo) Create(player model.Player) (err error) {
	exists, err := r.ExistsByID(player.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Player", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err := r.DB.Prepare(queryInsertPlayer)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(player)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}

// ResolveByID resolves a Player by its ID
func (r *PlayerMySQLRepo) ResolveByID(id uuid.UUID) (player *model.Player, err error) {
	player = &model.Player{}
	err = r.DB.Get(player, querySelectPlayer+" WHERE players.entity_id = ?", id)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	return
}

// TxUpdate transactionally updates a Player with the transaction object passed from elsewhere
func (r *PlayerMySQLRepo) TxUpdate(tx *sqlx.Tx, player model.Player) (err error) {
	stmt, err := tx.PrepareNamed(queryUpdatePlayer)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(player)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}
