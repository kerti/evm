package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/02-kitara-store/database"
	"github.com/kerti/evm/02-kitara-store/model"
	"github.com/kerti/evm/02-kitara-store/util/logger"
)

const (
	querySelectOrder = `
		SELECT
			orders.entity_id,
			orders.order_code,
			orders.total_price,
			orders.status
		FROM ` + "`orders`"

	querySelectOrderItem = `
		SELECT
			order_items.entity_id,
			order_items.order_entity_id,
			order_items.product_entity_id,
			order_items.qty,
			order_items.price
		FROM order_items`

	queryUpdateOrder = `
		UPDATE orders
		SET
			order_code = :order_code,
			total_price = :total_price,
			status = :status
		WHERE entity_id = :entity_id`
)

// Order is the Order repository interface
type Order interface {
	Startup()
	Shutdown()
	ResolveByID(id uuid.UUID) (order *model.Order, err error)
	TxUpdate(tx *sqlx.Tx, order model.Order) (err error)
}

// OrderMySQLRepo is the repository for Orders implemented with MySQL backend
type OrderMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup performs startup functions
func (r *OrderMySQLRepo) Startup() {
	logger.Trace("Order Repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *OrderMySQLRepo) Shutdown() {
	logger.Trace("Order Repository shutting down...")
}

// ResolveByID resolves an Order by its ID, including its items
func (r *OrderMySQLRepo) ResolveByID(id uuid.UUID) (order *model.Order, err error) {
	order = &model.Order{}
	err = r.DB.Get(order, querySelectOrder+" WHERE `orders`.entity_id = ?", id)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	query, args, err := r.DB.In(querySelectOrderItem+" WHERE order_items.order_entity_id = ?", order.ID)
	if err != nil {
		logger.Err("%v", err)
		return
	}

	orderItems := make([]model.OrderItem, 0)
	err = r.DB.Select(&orderItems, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	order.AttachItems(orderItems)

	return
}

// TxUpdate performs an update transactionally with the transaction object supplied from elsewhere
func (r *OrderMySQLRepo) TxUpdate(tx *sqlx.Tx, order model.Order) (err error) {
	stmt, err := tx.PrepareNamed(queryUpdateOrder)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(order)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
}
