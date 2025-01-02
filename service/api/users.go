package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) SearchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var request requests.Username

	request.Username = r.URL.Query().Get("username")

	if !request.IsValid() {
		utility.LogWithField("SearchUsers: invalid parameter - 'username'", http.StatusBadRequest, "username", request.Username, w, ctx)
		return
	}

	response, err := rt.db.SearchUsersByUsername(ctx.Uid, request.Username)

	if err != nil {
		utility.LogWithError("SearchUsers: error while searching users - SearchUsersByUsername", http.StatusInternalServerError, err, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (rt *_router) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToGet, requestErrors := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestErrors != nil {
		utility.LogWithError("GetUser: error while parsing request", http.StatusBadRequest, requestErrors, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToGet)

	if checkErrors != nil {
		utility.LogWithError("GetUser: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("GetUser: the requested user does not exists!", http.StatusNotFound, "userId", userToGet, w, ctx)
		return
	}

	isUserBlocker, checkErrors := rt.db.CheckBlock(ctx.Uid, userToGet)

	if checkErrors != nil {
		utility.LogWithError("GetUser: error while checking if the requested user is blocked - CheckBlock", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if isUserBlocker {
		utility.LogWithField("GetUser: the requester user is blocked!", http.StatusNotFound, "userId", userToGet, w, ctx)
		return
	}

	userRequested, requestError := rt.db.GetUser(ctx.Uid, userToGet)

	if requestError != nil {
		utility.LogWithError("GetUser: error while getting the requested user - GetUser", http.StatusInternalServerError, requestError, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(userRequested)
}

func (rt *_router) GetStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToGetStream, requestErrors := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestErrors != nil {
		utility.LogWithError("GetStream: error while parsing request", http.StatusBadRequest, requestErrors, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToGetStream)

	if checkErrors != nil {
		utility.LogWithError("GetStream: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("GetStream: the requested user does not exists!", http.StatusNotFound, "userId", userToGetStream, w, ctx)
		return
	}

	followedParam := r.URL.Query().Get("followed")
	var followedMode bool
	if followedParam == "" {
		followedMode = false
	} else {
		parsedFollowedMode, parseError := strconv.ParseBool(followedParam)
		if parseError != nil {
			followedMode = false
		} else {
			followedMode = parsedFollowedMode
		}
	}

	streamRequested, requestError := rt.db.GetPosts(ctx.Uid, userToGetStream, followedMode)

	if requestError != nil {
		utility.LogWithError("GetStream: error while getting the requested stream - GetPosts", http.StatusInternalServerError, requestError, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(streamRequested)
}
