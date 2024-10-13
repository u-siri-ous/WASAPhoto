package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) SetUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var request requests.Username

	requestErrors := json.NewDecoder(r.Body).Decode(&request)

	if requestErrors != nil {
		utility.LogWithError("SetUsername: error while parsing request", http.StatusBadRequest, requestErrors, w, ctx)
		return
	}

	if !request.IsValid() {
		utility.LogWithField("SetUsername: invalid param - 'username'", http.StatusBadRequest, "username", request.Username, w, ctx)
		return
	}

	isUsernameTaken, checkUsernameAvailabilityErrors := rt.db.CheckUsernameAvailability(request.Username)

	if requestErrors != nil {
		utility.LogWithError("SetUsername: error while checking the availability of the username - CheckUsernameAvailability", http.StatusInternalServerError, checkUsernameAvailabilityErrors, w, ctx)
		return
	}

	if isUsernameTaken {
		utility.LogWithField("SetUsername: the username is already taken", http.StatusForbidden, "username", request.Username, w, ctx)
		return
	}

	updateErrors := rt.db.SetUsername(ctx.Uid, request.Username)

	if updateErrors != nil {
		utility.LogWithError("SetUsername: error while updating the username - SetUsername", http.StatusInternalServerError, updateErrors, w, ctx)
		return
	}

}
