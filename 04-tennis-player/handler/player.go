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
	HandleCreate(w http.ResponseWriter, r *http.Request)
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

// HandleCreate handles the request
// @Summary Create a Player.
// @Description Creates a new Player.
// @Tags players
// @Accept json
// @Produce json
// @Param input body model.PlayerInput true "Input in the form of Player JSON."
// @Success 200 {object} response.BaseResponse{data=model.Player}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /players [post]
func (h *PlayerImpl) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input model.PlayerInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	player, err := h.PlayerService.Create(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, player)
}

// HandleAddBall handles the request
// @Summary Add balls.
// @Description Add balls to containers belonging to a particular user.
// @Tags players
// @Accept json
// @Produce json
// @Param input body model.PlayerAddBallInput true "Input specifying the player ID."
// @Success 200 {object} response.BaseResponse{data=model.Player}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /players/addBall [post]
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
