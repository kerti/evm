package service

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/02-kitara-store/database"
	"github.com/kerti/evm/02-kitara-store/model"
	"github.com/kerti/evm/02-kitara-store/repository"
	"github.com/kerti/evm/02-kitara-store/util/logger"
)

// Order is the service provider interface
type Order interface {
	Startup()
	Shutdown()
	Process(input model.OrderProcessInput) (*model.Order, error)
}

// OrderImpl is the service provider implementation
type OrderImpl struct {
	InventoryRepository repository.Inventory `inject:"inventoryRepository"`
	OrderRepository     repository.Order     `inject:"orderRepository"`
	DB                  *database.MySQL      `inject:"mysql"`
}

// Startup performs startup functions
func (s *OrderImpl) Startup() {
	logger.Trace("Order service starting up...")
}

// Shutdown cleans up everything and shuts down
func (s *OrderImpl) Shutdown() {
	logger.Trace("Order service shutting down...")
}

// Process processes an order
func (s *OrderImpl) Process(input model.OrderProcessInput) (*model.Order, error) {

	order, err := s.OrderRepository.ResolveByID(input.OrderID)
	if err != nil {
		return nil, err
	}

	productIDs := make([]uuid.UUID, 0)
	qtyMap := make(map[uuid.UUID]int)
	for _, orderItem := range order.Items {
		productIDs = append(productIDs, orderItem.ProductID)
		qtyMap[orderItem.ProductID] = orderItem.Qty
	}

	inventories, err := s.InventoryRepository.ResolveByProductIDs(productIDs)
	if err != nil {
		return nil, err
	}

	processedInventories := make([]model.Inventory, 0)
	for _, inventory := range inventories {
		err := inventory.Reserve(qtyMap[inventory.ProductID])
		if err != nil {
			return nil, err
		}
		processedInventories = append(processedInventories, inventory)
	}

	err = order.Process()
	if err != nil {
		return nil, err
	}

	err = s.DB.WithTransaction(s.DB, func(tx *sqlx.Tx, e chan error) {
		for _, inventory := range processedInventories {
			logger.Trace("updating inventory")
			if err := s.InventoryRepository.TxUpdate(tx, inventory); err != nil {
				e <- err
				return
			}
		}

		logger.Trace("updating order")
		if err := s.OrderRepository.TxUpdate(tx, *order); err != nil {
			e <- err
			return
		}

		e <- nil
	})

	return order, err

}
