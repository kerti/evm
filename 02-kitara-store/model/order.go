package model

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Order represents an Order entity
type Order struct {
	ID         uuid.UUID   `json:"id" db:"entity_id" validate:"min=36,max=36"`
	Code       string      `json:"code" db:"order_code"`
	TotalPrice float64     `json:"totalPrice" db:"total_price" validate:"min=0"`
	Status     string      `json:"status" db:"status"`
	Items      []OrderItem `json:"items" db:"-"`
}

// AttachItems attaches Order Items to an Order
func (o *Order) AttachItems(items []OrderItem) Order {
	for _, item := range items {
		if item.OrderID == o.ID {
			o.Items = append(o.Items, item)
		}
	}
	return *o
}

// Process updates an Order's status to processing
func (o *Order) Process() error {
	if o.Status != "new" {
		return errors.New("cannot process an order that is not new")
	}

	o.Status = "processing"

	return nil
}

// OrderItem represents an Order Item entity
type OrderItem struct {
	ID        uuid.UUID `json:"id" db:"entity_id" validate:"min=36,max=36"`
	OrderID   uuid.UUID `json:"orderId" db:"order_entity_id" validate:"min=36,max=36"`
	ProductID uuid.UUID `json:"productId" db:"product_entity_id" validate:"min=36,max=36"`
	Qty       int       `json:"qty" db:"qty" validate:"min=1"`
	Price     float64   `json:"price" db:"price"`
}

// OrderProcessInput represents an input where the user wants to process an Order
type OrderProcessInput struct {
	OrderID uuid.UUID `json:"orderId"`
}
