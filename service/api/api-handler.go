package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login
	rt.router.POST("/login", rt.wrap(rt.Login, false))
	// Search Users
	rt.router.GET("/users/", rt.wrap(rt.SearchUsers, true))
	// Get User
	rt.router.GET("/users/:userId", rt.wrap(rt.GetUser, true))
	// Settings
	rt.router.PUT("/settings/username", rt.wrap(rt.SetUsername, true))
	// Block
	rt.router.PUT("/blocked/:userId", rt.wrap(rt.BlockUser, true))
	rt.router.DELETE("/blocked/:userId", rt.wrap(rt.UnblockUser, true))
	rt.router.GET("/blocked/", rt.wrap(rt.BlockedUsers, true))
	// Follow
	rt.router.GET("/followers/:userId", rt.wrap(rt.GetFollowers, true))
	rt.router.GET("/following/:userId", rt.wrap(rt.GetFollowing, true))
	rt.router.PUT("/following/:userId", rt.wrap(rt.FollowUser, true))
	rt.router.DELETE("/following/:userId", rt.wrap(rt.UnfollowUser, true))
	// Post
	rt.router.POST("/posts/", rt.wrap(rt.CreatePost, true))
	rt.router.DELETE("/posts/:postId", rt.wrap(rt.DeletePost, true))
	rt.router.GET("/posts/:postId/likes/self", rt.wrap(rt.Likes, true))
	rt.router.POST("/posts/:postId/likes/self", rt.wrap(rt.LikePost, true))
	rt.router.DELETE("/posts/:postId/likes/self", rt.wrap(rt.UnlikePost, true))
	rt.router.GET("/posts/:postId/comments/", rt.wrap(rt.Comments, true))
	rt.router.POST("/posts/:postId/comments/", rt.wrap(rt.CommentPost, true))
	rt.router.DELETE("/posts/:postId/comments/:commentId", rt.wrap(rt.DeleteCommentPost, true))

	// Liveness check
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
