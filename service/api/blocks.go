package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
)

func (rt *_router) BlockUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToBan, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("BlockUser: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToBan)

	if checkErrors != nil {
		utility.LogWithError("BlockUser: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("BlockUser: the requested user does not exists!", http.StatusNotFound, "userId", userToBan, w, ctx)
		return
	}

	userAlreadyBanned, checkAlreadyBannedErrors := rt.db.CheckBlock(userToBan, ctx.Uid)

	if checkAlreadyBannedErrors != nil {
		utility.LogWithError("BlockUser: error while checking the block blocking the user - CheckBlock", http.StatusInternalServerError, checkAlreadyBannedErrors, w, ctx)
		return
	}

	if userAlreadyBanned {
		utility.LogWithField("BlockUser: the requester already blocked the user", http.StatusBadRequest, "userId", userToBan, w, ctx)
		return
	}

	if ctx.Uid == userToBan {
		utility.LogWithField("BlockUser: the requester and the user to ban cannot be the same", http.StatusBadRequest, "userId", userToBan, w, ctx)
		return
	}

	insertError := rt.db.BlockUser(ctx.Uid, userToBan)

	if insertError != nil {
		utility.LogWithError("BlockUser: error while blocking the user - BlockUser", http.StatusInternalServerError, insertError, w, ctx)
		return
	}

	unfollowError := rt.db.UnfollowUser(ctx.Uid, userToBan)

	if unfollowError != nil {
		utility.LogWithError("BlockUser: error while unfollowing the user after block - BlockUser", http.StatusInternalServerError, unfollowError, w, ctx)
		return
	}

	unfollowReversedError := rt.db.UnfollowUser(userToBan, ctx.Uid)

	if unfollowReversedError != nil {
		utility.LogWithError("BlockUser: error while unfollowing the blocker - BlockUser", http.StatusInternalServerError, unfollowReversedError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) UnblockUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userToUnban, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("UnblockUser: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(userToUnban)

	if checkErrors != nil {
		utility.LogWithError("UnblockUser: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("UnblockUser: the requested user does not exists!", http.StatusNotFound, "userId", userToUnban, w, ctx)
		return
	}

	if ctx.Uid == userToUnban {
		utility.LogWithField("UnblockUser: the requester and the user to unban cannot be the same", http.StatusBadRequest, "userId", userToUnban, w, ctx)
		return
	}

	insertError := rt.db.UnblockUser(ctx.Uid, userToUnban)

	if insertError != nil {
		utility.LogWithError("UnblockUser: error while unblocking the user - UnblockUser", http.StatusInternalServerError, insertError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) BlockedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	blockedUsers, requestErrors := rt.db.ListOfBlocker(ctx.Uid)

	if requestErrors != nil {
		utility.LogWithError("BlockedUsers: error while retrieving the blocked users - ListOfBlocker", http.StatusInternalServerError, requestErrors, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(blockedUsers)
}
