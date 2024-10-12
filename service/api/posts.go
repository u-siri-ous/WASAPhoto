package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
	"github.com/u-siri-ous/WASAPhoto/service/api/reqcontext"
	"github.com/u-siri-ous/WASAPhoto/service/globaltime"
	"github.com/u-siri-ous/WASAPhoto/service/storage"
)

func (rt *_router) CreatePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestError := r.ParseMultipartForm(storage.MaxRequestFileSize)
	if requestError != nil {
		utility.LogWithError("CreatePost: error while parsing request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	photoCaption := r.FormValue("caption")
	photo, photoHeader, getFileError := r.FormFile("photo")
	if getFileError != nil {
		utility.LogWithError("CreatePost: error while getting the file to upload", http.StatusBadRequest, getFileError, w, ctx)
		return
	}
	defer photo.Close()

	if photoHeader.Header.Get("Content-Type") != storage.AllowedMimeType {
		utility.LogError("CreatePost: unsupported file", http.StatusBadRequest, w, ctx)
		return
	}

	if photoHeader.Size > storage.MaxFileSize {
		utility.LogError("CreatePost: the file is too large", http.StatusRequestEntityTooLarge, w, ctx)
		return
	}

	now := globaltime.Now()
	postId, insertError := rt.db.CreatePost(ctx.Uid, photoCaption, now)

	if insertError != nil {
		utility.LogWithError("CreatePost: error while creating the post", http.StatusInternalServerError, insertError, w, ctx)
		return
	}

	uploadError := storage.SavePhoto(photo, postId, ctx.Uid)
	if uploadError != nil {
		rollbackError := rt.db.DeletePost(ctx.Uid, postId)
		if rollbackError != nil {
			utility.LogWithError("CreatePost: error while executing a post rollback", http.StatusInternalServerError, rollbackError, w, ctx)
			return
		}
		utility.LogWithError("CreatePost: error while saving the photo", http.StatusInternalServerError, uploadError, w, ctx)
		return
	}

	post := structs.Post{
		Id:             postId,
		Author:         ctx.Uid,
		Caption:        photoCaption,
		Likes:          0,
		Comments:       0,
		TimeOfCreation: now,
		IsLiked:        false,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(post)
}

func (rt *_router) DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postToDelete, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)
	if requestError != nil {
		utility.LogWithError("DeletePost: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToDelete)
	if checkPostError != nil {
		utility.LogWithError("DeletePost: error while checking the request", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("DeletePost: the requested post does not exists", http.StatusBadRequest, "postId", postToDelete, w, ctx)
		return
	}

	deletePostError := rt.db.DeletePost(ctx.Uid, postToDelete)
	if deletePostError != nil {
		utility.LogWithError("DeletePost: error while deleting the post", http.StatusInternalServerError, deletePostError, w, ctx)
		return
	} else {
		deletePhotoError := storage.DeletePhoto(ctx.Uid, postToDelete)
		if deletePhotoError != nil {
			utility.LogWithError("DeletePost: error while deleting the photo from the storage", http.StatusInternalServerError, deletePostError, w, ctx)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
