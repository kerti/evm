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

// Container is the handler interface for Containers
type Container interface {
	Startup()
	Shutdown()
	HandleCreate(w http.ResponseWriter, r *http.Request)
}

// ContainerImpl is the handler implementation for Containers
type ContainerImpl struct {
	ContainerService service.Container `inject:"containerService"`
}

// Startup performs startup functions
func (h *ContainerImpl) Startup() {
	logger.Trace("Container Handler starting up...")
}

// Shutdown cleans up everything and shuts down
func (h *ContainerImpl) Shutdown() {
	logger.Trace("Container Handler shutting down...")
}

// HandleCreate handles the request
// @Summary Create a Container.
// @Description Creates a new Container.
// @Tags containers
// @Accept json
// @Produce json
// @Param input body model.ContainerInput true "Input in the form of Container JSON."
// @Success 200 {object} response.BaseResponse{data=model.Container}
// @Failure 400 {object} response.BaseResponse
// @Failure 409 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /containers [post]
func (h *ContainerImpl) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input model.ContainerInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
		return
	}

	container, err := h.ContainerService.Create(input)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusCreated, container)
}
