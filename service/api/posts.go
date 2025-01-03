package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/u-siri-ous/WASAPhoto/service/api/components/requests"
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
		utility.LogError("CreatePost: unsupported file: "+photoHeader.Header.Get("Content-Type"), http.StatusUnsupportedMediaType, w, ctx)
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
		AuthorId:       ctx.Uid,
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

func (rt *_router) LikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postToLike, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("LikePost: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToLike)
	if checkPostError != nil {
		utility.LogWithError("LikePost: error while checking the request", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("LikePost: the requested post does not exists", http.StatusBadRequest, "postId", postToLike, w, ctx)
		return
	}

	alreadyLiked, checkLikeError := rt.db.CheckLikePost(ctx.Uid, postToLike)
	if checkLikeError != nil {
		utility.LogWithError("LikePost: error while checking if the post is already liked", http.StatusInternalServerError, checkLikeError, w, ctx)
		return
	}

	if alreadyLiked > 0 {
		utility.LogWithField("LikePost: the post is already liked!", http.StatusBadRequest, "postId", postToLike, w, ctx)
		return
	}

	insertLikeError := rt.db.InsertLikePost(ctx.Uid, postToLike)

	if insertLikeError != nil {
		utility.LogWithError("LikePost: error during the like process", http.StatusBadRequest, insertLikeError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) UnlikePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postToLike, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("UnlikePost: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToLike)
	if checkPostError != nil {
		utility.LogWithError("UnlikePost: error while checking the request", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("UnlikePost: the requested post does not exists", http.StatusBadRequest, "postId", postToLike, w, ctx)
		return
	}

	alreadyLiked, checkLikeError := rt.db.CheckLikePost(ctx.Uid, postToLike)
	if checkLikeError != nil {
		utility.LogWithError("UnlikePost: error while checking if the post is already liked", http.StatusInternalServerError, checkLikeError, w, ctx)
		return
	}

	if alreadyLiked <= 0 {
		utility.LogWithField("UnlikePost: the post is not liked!", http.StatusBadRequest, "postId", postToLike, w, ctx)
		return
	}

	insertLikeError := rt.db.DeleteLikePost(ctx.Uid, postToLike)

	if insertLikeError != nil {
		utility.LogWithError("UnlikePost: error during the unlike process", http.StatusBadRequest, insertLikeError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) Likes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postToCheckLike, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("Likes: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToCheckLike)
	if checkPostError != nil {
		utility.LogWithError("Likes: error while checking the request", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("Likes: the requested post does not exists", http.StatusNotFound, "postId", postToCheckLike, w, ctx)
		return
	}

	userListResponse, getLikesError := rt.db.GetLikes(ctx.Uid, postToCheckLike)
	if getLikesError != nil {
		utility.LogWithError("Likes: error while getting the list of likes - GetLikes", http.StatusInternalServerError, getLikesError, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(userListResponse)
}

func (rt *_router) CommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var request requests.Comment

	postToComment, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("CommentPost: error while parsing the request - path", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	requestBodyError := json.NewDecoder(r.Body).Decode(&request)

	if requestBodyError != nil {
		utility.LogWithError("CommentPost: error while parsing the request - body", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	if !request.IsValid() {
		utility.LogWithField("CommentPost: error while validating the request - body", http.StatusBadRequest, "text", request.Text, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToComment)
	if checkPostError != nil {
		utility.LogWithError("CommentPost: error while checking the request - CheckPostId", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("CommentPost: the requested post does not exists", http.StatusNotFound, "postId", postToComment, w, ctx)
		return
	}

	insertError := rt.db.InsertCommentPost(ctx.Uid, postToComment, request.Text)

	if insertError != nil {
		utility.LogWithError("CommentPost: error while creating the comment - InsertCommentPost", http.StatusBadRequest, insertError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) Comments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postToGetComments, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("Comments: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postToGetComments)
	if checkPostError != nil {
		utility.LogWithError("Comments: error while checking the request", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("Comments: the requested post does not exists", http.StatusNotFound, "postId", postToGetComments, w, ctx)
		return
	}

	postCommentsResponse, getCommentsError := rt.db.GetCommentsPost(ctx.Uid, postToGetComments)
	if getCommentsError != nil {
		utility.LogWithError("Comments: error while getting the list of comments - GetCommentsPost", http.StatusInternalServerError, getCommentsError, w, ctx)
		return
	}

	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(postCommentsResponse)
}

func (rt *_router) DeleteCommentPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postId, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("CommentPost: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	commentToDelete, requestError := strconv.ParseUint(ps.ByName("commentId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("CommentPost: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postId)
	if checkPostError != nil {
		utility.LogWithError("CommentPost: error while checking the request - CheckPostId", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("CommentPost: the requested post does not exists", http.StatusNotFound, "postId", postId, w, ctx)
		return
	}

	commentExists, checkCommentError := rt.db.CheckCommentId(commentToDelete)

	if checkCommentError != nil {
		utility.LogWithError("CommentPost: error while checking the request - CheckCommentId", http.StatusInternalServerError, checkCommentError, w, ctx)
		return
	}

	if !commentExists {
		utility.LogWithField("CommentPost: the requested post does not exists", http.StatusNotFound, "commentId", commentToDelete, w, ctx)
		return
	}

	deleteError := rt.db.DeleteCommentPost(ctx.Uid, postId, commentToDelete)

	if deleteError != nil {
		utility.LogWithError("CommentPost: error while deleting the comment - DeleteCommentPost", http.StatusBadRequest, deleteError, w, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) GetPostPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postId, requestError := strconv.ParseUint(ps.ByName("postId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("GetPostPhoto: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	posterUserId, requestError := strconv.ParseUint(ps.ByName("userId"), 10, 64)

	if requestError != nil {
		utility.LogWithError("GetPostPhoto: error while parsing the request", http.StatusBadRequest, requestError, w, ctx)
		return
	}

	userExists, checkErrors := rt.db.CheckUserId(posterUserId)

	if checkErrors != nil {
		utility.LogWithError("GetPostPhoto: error while checking the requested user on the database - CheckUserId", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if !userExists {
		utility.LogWithField("GetPostPhoto: the requested user does not exists!", http.StatusNotFound, "userId", posterUserId, w, ctx)
		return
	}

	isUserBlocked, checkErrors := rt.db.CheckBlock(posterUserId, ctx.Uid)

	if checkErrors != nil {
		utility.LogWithError("GetPostPhoto: error while checking if the requested user is blocked - CheckBlock", http.StatusInternalServerError, checkErrors, w, ctx)
		return
	}

	if isUserBlocked {
		utility.LogWithField("GetPostPhoto: the requester user is blocked!", http.StatusNotFound, "userId", posterUserId, w, ctx)
		return
	}

	postExists, checkPostError := rt.db.CheckPostId(postId)
	if checkPostError != nil {
		utility.LogWithError("GetPostPhoto: error while checking the request - CheckPostId", http.StatusInternalServerError, checkPostError, w, ctx)
		return
	}

	if !postExists {
		utility.LogWithField("GetPostPhoto: the requested post does not exists", http.StatusNotFound, "postId", postId, w, ctx)
		return
	}

	path := filepath.Join(storage.BasePath, "/", strconv.FormatUint(posterUserId, 10), storage.UploadedPhotoFolder, fmt.Sprintf("%d.jpg", postId))
	photofile, openPhotoError := os.Open(path)
	if openPhotoError != nil {
		utility.LogWithError("GetPostPhoto: resource not found", http.StatusInternalServerError, openPhotoError, w, ctx)
		return
	} else {
		w.Header().Set("Content-Type", storage.AllowedMimeType)
		buf := bytes.NewBuffer(nil)
		_, copyPhotoError := io.Copy(buf, photofile)
		if copyPhotoError != nil {
			utility.LogWithError("GetPostPhoto: error while copying the file to the buffer", http.StatusInternalServerError, copyPhotoError, w, ctx)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			_, writeResponseError := w.Write(buf.Bytes())
			if writeResponseError != nil {
				utility.LogWithError("GetPostPhoto: error while writing the response", http.StatusInternalServerError, writeResponseError, w, ctx)
				return
			}
			return
		}
	}
}
