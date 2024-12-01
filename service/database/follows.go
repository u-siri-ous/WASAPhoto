package database

import "github.com/u-siri-ous/WASAPhoto/service/api/components/structs"

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

func (db *appdbimpl) GetFollowers(currentUserId uint64, userId uint64) (structs.UserList, error) {
	var result structs.UserList

	const getFollowersList = "SELECT follow.id, follow.username FROM ( SELECT u.id, u.username FROM users u LEFT JOIN follows f ON f.followerUserId = u.id LEFT JOIN users uFollower ON uFollower.id = f.followedUserId WHERE uFollower.id = ? ) follow left JOIN blocks b ON b.blockerUserId = follow.id WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL"

	rows, errors := db.c.Query(getFollowersList, userId, currentUserId)

	if errors != nil {
		return result, errors
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.UserReduced
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return result, err
		}
		result.Users = append(result.Users, user)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, errors
}

func (db *appdbimpl) GetFollowing(currentUserId uint64, userId uint64) (structs.UserList, error) {
	var result structs.UserList

	const getFollowingList = "SELECT follow.id, follow.username FROM ( SELECT u.id, u.username FROM users u LEFT JOIN follows f ON f.followedUserId = u.id LEFT JOIN users uFollowed ON uFollowed.id = f.followerUserId WHERE uFollowed.id = ? ) follow left JOIN blocks b ON b.blockerUserId = follow.id WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL"

	rows, errors := db.c.Query(getFollowingList, userId, currentUserId)

	if errors != nil {
		return result, errors
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.UserReduced
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return result, err
		}
		result.Users = append(result.Users, user)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, errors
}
