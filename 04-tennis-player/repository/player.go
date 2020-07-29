package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

const (
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
