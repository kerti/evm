package handler

import (
	"net/http"

	"github.com/kerti/evm/01-machine-gun/handler/response"
	"github.com/kerti/evm/01-machine-gun/util/logger"
)

// Health is the handler interface for Server Health
type Health interface {
	Startup()
	Shutdown()
	HandleHealthCheck(w http.ResponseWriter, r *http.Request)
}

// HealthImpl is the handler implementation for Server Health
type HealthImpl struct {
	isPreparingShutdown bool
	isHealthy           bool
}

// Startup perform startup functions
func (h *HealthImpl) Startup() {
	logger.Trace("Health Handler starting up...")
	h.isPreparingShutdown = false
	h.isHealthy = true
}

// PrepareShutdown prepares the service for shutdown
func (h *HealthImpl) PrepareShutdown() {
	logger.Trace("Health Handler preparing for shutdown...")
	h.isPreparingShutdown = true
	h.isHealthy = false
}

// Shutdown cleans up everything and shuts down
func (h *HealthImpl) Shutdown() {
	h.isPreparingShutdown = false
	logger.Trace("Health Handler shutting down...")
}

// HandleHealthCheck handles the request
func (h *HealthImpl) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if h.isHealthy {
		response.RespondWithMessage(w, http.StatusOK, "OK")
	} else {
		if h.isPreparingShutdown {
			response.RespondWithPreparingShutdown(w)
		} else {
			response.RespondWithUnhealthy(w)
		}
	}
}
