package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kerti/evm/04-tennis-player/handler/response"

	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/service"
	"github.com/kerti/evm/04-tennis-player/util/failure"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

// Player is the handler interface for Players
type Player interface {
	Startup()
	Shutdown()
	HandleAddBall(w http.ResponseWriter, r *http.Request)
}

// PlayerImpl is the handler implementation for Players
type PlayerImpl struct {
	PlayerService service.Player `inject:"playerService"`
}

// Startup performs startup functions
func (h *PlayerImpl) Startup() {
	logger.Trace("Player Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *PlayerImpl) Shutdown() {
	logger.Trace("Player Handler shutting down...")
}

// HandleAddBall handles the request
func (h *PlayerImpl) HandleAddBall(w http.ResponseWriter, r *http.Request) {
	var input model.PlayerAddBallInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	player, err := h.PlayerService.AddBall(input.PlayerID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, player)
}