package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req requests.Username

	requestErrors := json.NewDecoder(r.Body).Decode(&req)

	if requestErrors != nil {
		utility.LogWithError("Login: error while parsing request", http.StatusBadRequest, requestErrors, w, ctx)
		return
	}

	if !req.IsValid() {
		utility.LogWithField("Login: invalid param - 'username'", http.StatusBadRequest, "username", req.Username, w, ctx)
		return
	}

	uid, created, responseErrors := rt.db.InsertUser(req.Username)

	if responseErrors != nil {
		utility.LogWithError("Login: error in DB -'Insert_user(request.Username)'s ", http.StatusInternalServerError, responseErrors, w, ctx)
		return
	}

	ctx.Logger.Debug("Login: OK")

	result := structs.User{Id: uid, Username: req.Username}
	w.Header().Set("content-type", "application/json")
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	_ = json.NewEncoder(w).Encode(result)
}
