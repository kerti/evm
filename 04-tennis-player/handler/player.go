package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	HandleResolveByID(w http.ResponseWriter, r *http.Request)
	HandleResolvePage(w http.ResponseWriter, r *http.Request)
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

// HandleResolveByID handles the request
// @Summary Resolve a Player.
// @Description Resolves a Player by its ID.
// @Tags players
// @Produce json
// @Param id path string true "The Player's identifier."
// @Success 200 {object} response.BaseResponse{data=model.Player}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /players/{id} [get]
func (h *PlayerImpl) HandleResolveByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	player, err := h.PlayerService.ResolveByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, player)
}

// HandleResolvePage handles the request
// @Summary Resolve a Page of Players.
// @Description Resolves a Page of Players based on page and page size parameters.
// @Tags players
// @Produce json
// @Param page query int false "The page number. Defaults to 1."
// @Param pageSize query int false "The number of records on a page. Defaults to 10."
// @Success 200 {object} response.BaseResponse{data=model.Page}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /players/ [get]
func (h *PlayerImpl) HandleResolvePage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	pageNumStr, withPageNum := r.URL.Query()["page"]
	pageSizeStr, withPageSize := r.URL.Query()["pageSize"]

	pageNum := 1
	if withPageNum && len(pageNumStr[0]) > 0 {
		pageNum, err = strconv.Atoi(pageNumStr[0])
		if err != nil {
			response.RespondWithError(w, failure.BadRequest(err))
			return
		}
	}

	if pageNum <= 0 {
		response.RespondWithError(
			w,
			failure.BadRequestFromString("page must be positive integer"))
		return
	}

	pageSize := 10
	if withPageSize && len(pageSizeStr[0]) > 0 {
		pageSize, err = strconv.Atoi(pageSizeStr[0])
		if err != nil {
			response.RespondWithError(w, failure.BadRequest(err))
			return
		}
	}

	if pageSize <= 0 {
		response.RespondWithError(
			w,
			failure.BadRequestFromString("page size must be positive integer"))
		return
	}

	page, err := h.PlayerService.ResolvePage(pageNum, pageSize)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, page)
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
