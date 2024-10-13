package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login
	rt.router.POST("/login", rt.wrap(rt.Login, false))
	// Search Users
	rt.router.GET("/users", rt.wrap(rt.SearchUsers, true))
	// Get User
	rt.router.GET("/users/:userId", rt.wrap(rt.GetUser, true))
	// Settings
	rt.router.PUT("/settings/username", rt.wrap(rt.SetUsername, true))
	// Block
	rt.router.PUT("/block/:userId", rt.wrap(rt.BlockUser, true))
	rt.router.DELETE("/block/:userId", rt.wrap(rt.UnblockUser, true))
	rt.router.GET("/blocked", rt.wrap(rt.BlockedUsers, true))
	// Follow
	rt.router.PUT("/follow/:userId", rt.wrap(rt.FollowUser, true))
	rt.router.DELETE("/follow/:userId", rt.wrap(rt.UnfollowUser, true))
	rt.router.GET("/followers/:userId", rt.wrap(rt.GetFollowers, true))
	rt.router.GET("/following/:userId", rt.wrap(rt.GetFollowing, true))
	// Post
	rt.router.POST("/posts", rt.wrap(rt.CreatePost, true))
	rt.router.DELETE("/posts/:postId", rt.wrap(rt.DeletePost, true))
	rt.router.GET("/posts/:postId/likes", rt.wrap(rt.Likes, true))
	rt.router.POST("/posts/:postId/likes", rt.wrap(rt.LikePost, true))
	rt.router.DELETE("/posts/:postId/likes", rt.wrap(rt.UnlikePost, true))
	rt.router.POST("/posts/:postId/comments", rt.wrap(rt.CommentPost, true))
	rt.router.DELETE("/posts/:postId/comments/:commentId", rt.wrap(rt.DeleteCommentPost, true))

	return rt.router
}
