package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/kerti/evm/04-tennis-player/handler/response"
	"github.com/kerti/evm/04-tennis-player/util/failure"
)

func getIDFromRequest(w http.ResponseWriter, r *http.Request) (id uuid.UUID, err error) {
	vars := mux.Vars(r)

	id, err = uuid.FromString(vars["id"])
	if err != nil {
		response.RespondWithError(w, failure.BadRequest(err))
	}

	return
}
