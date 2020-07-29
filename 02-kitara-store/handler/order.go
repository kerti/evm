package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kerti/evm/02-kitara-store/handler/response"
	"github.com/kerti/evm/02-kitara-store/model"

	"github.com/kerti/evm/02-kitara-store/service"
	"github.com/kerti/evm/02-kitara-store/util/failure"
	"github.com/kerti/evm/02-kitara-store/util/logger"
)

// Order is the handler interface for Orders
type Order interface {
	Startup()
	Shutdown()
	HandleProcessOrder(w http.ResponseWriter, r *http.Request)
}

// OrderImpl is the handler implementation for Orders
type OrderImpl struct {
	Service service.Order `inject:"orderService"`
}

// Startup performs startup functions
func (h *OrderImpl) Startup() {
	logger.Trace("Order Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *OrderImpl) Shutdown() {
	logger.Trace("Order Handler shutting down...")
}

// HandleProcessOrder handles the request
func (h *OrderImpl) HandleProcessOrder(w http.ResponseWriter, r *http.Request) {
	var input model.OrderProcessInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	order, err := h.Service.Process(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, order)
}
