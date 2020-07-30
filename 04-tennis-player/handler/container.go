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

// Container is the handler interface for Containers
type Container interface {
	Startup()
	Shutdown()
	HandleResolveByID(w http.ResponseWriter, r *http.Request)
	HandleResolvePage(w http.ResponseWriter, r *http.Request)
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

// HandleResolveByID handles the request
// @Summary Resolve a Container.
// @Description Resolves a Container by its ID.
// @Tags containers
// @Produce json
// @Param id path string true "The Container's identifier."
// @Success 200 {object} response.BaseResponse{data=model.Container}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /containers/{id} [get]
func (h *ContainerImpl) HandleResolveByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(w, r)
	if err != nil {
		return
	}

	container, err := h.ContainerService.ResolveByID(id)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, container)
}

// HandleResolvePage handles the request
// @Summary Resolve a Page of Containers.
// @Description Resolves a Page of Containers based on page and page size parameters.
// @Tags containers
// @Produce json
// @Param page query int false "The page number. Defaults to 1."
// @Param pageSize query int false "The number of records on a page. Defaults to 10."
// @Success 200 {object} response.BaseResponse{data=model.Page}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /containers/ [get]
func (h *ContainerImpl) HandleResolvePage(w http.ResponseWriter, r *http.Request) {
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

	page, err := h.ContainerService.ResolvePage(pageNum, pageSize)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	response.RespondWithJSON(w, http.StatusOK, page)
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
