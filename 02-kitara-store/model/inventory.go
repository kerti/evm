package model

import (
	"errors"

	"github.com/gofrs/uuid"
)

// Inventory represents an Inventory object
type Inventory struct {
	ID           uuid.UUID `json:"id" db:"entity_id" validate:"min=36,max=36"`
	ProductID    uuid.UUID `json:"productId" db:"product_entity_id" validate:"min=36,max=36"`
	QtyInStore   int       `json:"qtyInStore" db:"qty_in_store" validate:"min=0"`
	QtyReserved  int       `json:"qtyReserved" db:"qty_reserved" validate:"min=0"`
	QtyAvailable int       `json:"qtyAvailable" db:"qty_available" validate:"min=0"`
}

// Reserve reserves the specified amount of inventory
func (i *Inventory) Reserve(qty int) error {
	i.QtyReserved += qty
	i.QtyAvailable = i.QtyInStore - i.QtyReserved

	return i.Validate()
}

// Validate validates the Inventory object
func (i *Inventory) Validate() error {
	if i.QtyInStore < 0 {
		return errors.New("cannot have negative in-store quantity")
	}

	if i.QtyReserved > i.QtyInStore {
		return errors.New("cannot reserve more than in-store quantity")
	}

	if i.QtyAvailable < 0 {
		return errors.New("cannot have negative available quantity")
	}

	return nil
}
