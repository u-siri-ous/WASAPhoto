package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) FollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToFollow, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("FollowUser: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToFollow)

	if checkErrors != nil {
		utility.LogWithError("FollowUser: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("FollowUser: the requested user does not exists!", http.StatusNotFound, "userId", userToFollow, w, ctx)
		return
	}

	insertErrors := rt.db.FollowUser(ctx.Uid, userToFollow)

	if insertErrors != nil {
		utility.LogWithError("FollowUser: error while following the user - FollowUser", http.StatusInternalServerError, insertErrors, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToFollow, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("UnfollowUser: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToFollow)

	if checkErrors != nil {
		utility.LogWithError("UnfollowUser: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("UnfollowUser: the requested user does not exists!", http.StatusNotFound, "userId", userToFollow, w, ctx)
		return
	}

	deleteErrors := rt.db.UnfollowUser(ctx.Uid, userToFollow)

	if deleteErrors != nil {
		utility.LogWithError("UnfollowUser: error while unfollowing the user - UnfollowUser", http.StatusInternalServerError, deleteErrors, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) GetFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("GetFollowers: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userId)

	if checkErrors != nil {
		utility.LogWithError("GetFollowers: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("GetFollowers: the requested user does not exists!", http.StatusNotFound, "userId", userId, w, ctx)
		return
	}

	isUserBlocker, checkErrors := rt.db.CheckBlock(ctx.Uid, userId)

	if checkErrors != nil {
		utility.LogWithError("GetFollowers: error while checking if the requested user is blocked - CheckBlock", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if isUserBlocker {
		utility.LogWithField("GetFollowers: the requester user is blocked!", http.StatusForbidden, "userId", userId, w, ctx)
		return
	}

	result, errors := rt.db.GetFollowers(ctx.Uid, userId)

	if errors != nil {
		utility.LogWithError("GetFollowers: error while getting the followers - GetFollowers", http.StatusInternalServerError, errors, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func (rt *_router) GetFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("GetFollowing: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userId)

	if checkErrors != nil {
		utility.LogWithError("GetFollowing: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("GetFollowing: the requested user does not exists!", http.StatusNotFound, "userId", userId, w, ctx)
		return
	}

	isUserBlocker, checkErrors := rt.db.CheckBlock(ctx.Uid, userId)

	if checkErrors != nil {
		utility.LogWithError("GetFollowing: error while checking if the requested user is blocked - CheckBlock", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if isUserBlocker {
		utility.LogWithField("GetFollowing: the requester user is blocked!", http.StatusForbidden, "userId", userId, w, ctx)
		return
	}

	result, errors := rt.db.GetFollowing(ctx.Uid, userId)

	if errors != nil {
		utility.LogWithError("GetFollowing: error while getting the followed - GetFollowing", http.StatusInternalServerError, errors, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
