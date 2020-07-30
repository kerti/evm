package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/02-kitara-store/database"
	"github.com/kerti/evm/02-kitara-store/model"
	"github.com/kerti/evm/02-kitara-store/util/logger"
)

const (
	querySelectInventory = `
		SELECT
			inventory.entity_id,
			inventory.product_entity_id,
			inventory.qty_in_store,
			inventory.qty_reserved,
			inventory.qty_available
		FROM inventory`

	queryUpdateInventory = `
		UPDATE inventory
		SET
			product_entity_id = :product_entity_id,
			qty_in_store = :qty_in_store,
			qty_reserved = :qty_reserved,
			qty_available = :qty_available
		WHERE entity_id = :entity_id`
)

// Inventory is the Inventory repository interface
type Inventory interface {
	Startup()
	Shutdown()
	ResolveByProductIDs(ids []uuid.UUID) (inventories []model.Inventory, err error)
	TxUpdate(tx *sqlx.Tx, inventory model.Inventory) (err error)
}

// InventoryMySQLRepo is the repository for Inventory implemented with MySQL backend
type InventoryMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup performs startup functions
func (r *InventoryMySQLRepo) Startup() {
	logger.Trace("Inventory Repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *InventoryMySQLRepo) Shutdown() {
	logger.Trace("Inventory Repository shutting down...")
}

// ResolveByProductIDs resolves Inventories by their Product IDs
func (r *InventoryMySQLRepo) ResolveByProductIDs(ids []uuid.UUID) (inventories []model.Inventory, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(querySelectInventory+" WHERE inventory.product_entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&inventories, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// TxUpdate performs an update transactionally with transaction object supplied from elsewhere
func (r *InventoryMySQLRepo) TxUpdate(tx *sqlx.Tx, inventory model.Inventory) (err error) {
	stmt, err := tx.PrepareNamed(queryUpdateInventory)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(inventory)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}
