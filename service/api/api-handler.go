package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	// Login
	rt.router.POST("/login", rt.wrap(rt.Login, false))
	// Search Users
	rt.router.GET("/users", rt.wrap(rt.SearchUsers, true))
	// Get User
	rt.router.GET("/users/:userId", rt.wrap(rt.GetUser, true))
	// TODO: update username (tag settings)
	// Block
	rt.router.PUT("/block/:userId", rt.wrap(rt.BlockUser, true))
	rt.router.DELETE("/block/:userId", rt.wrap(rt.UnblockUser, true))
	rt.router.GET("/blocked", rt.wrap(rt.BlockedUsers, true))
	// Follow
	rt.router.PUT("/follow/:userId", rt.wrap(rt.FollowUser, true))
	rt.router.DELETE("/follow/:userId", rt.wrap(rt.UnfollowUser, true))
	// TODO: lista followers e following
	// Post
	rt.router.POST("/posts", rt.wrap(rt.CreatePost, true))
	rt.router.DELETE("/posts/:postId", rt.wrap(rt.DeletePost, true))
	return rt.router
}
