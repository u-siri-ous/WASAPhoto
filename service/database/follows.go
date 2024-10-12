package database

func (db *appdbimpl) FollowUser(currentUserId uint64, userToFollow uint64) error {
	const followUserQuery = "INSERT INTO follows(followerUserId, followedUserId) VALUES (?, ?)"

	_, responseErrors := db.c.Exec(followUserQuery, currentUserId, userToFollow)

	return responseErrors
}

func (db *appdbimpl) UnfollowUser(currentUserId uint64, userToUnfollow uint64) error {
	const unfollowUserQuery = "DELETE FROM follows WHERE followerUserId = ? AND followedUserId = ?"

	_, responseErrors := db.c.Exec(unfollowUserQuery, currentUserId, userToUnfollow)

	return responseErrors
}
