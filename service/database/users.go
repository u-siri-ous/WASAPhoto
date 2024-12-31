package database

import (
	"database/sql"
	"errors"

	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
)

func (db *appdbimpl) CheckUserId(userId uint64) (bool, error) {
	var response bool
	const checkUserQuery = `SELECT EXISTS (SELECT 1 FROM users WHERE id = ?)`
	responseError := db.c.QueryRow(checkUserQuery, userId).Scan(&response)

	return response, responseError
}

func (db *appdbimpl) GetUserIdByUsername(username string) (uint64, error) {
	var id uint64
	err := db.c.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&id)
	return id, err
}

func (db *appdbimpl) InsertUser(username string) (uint64, bool, error) {
	id, err := db.GetUserIdByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {

		result, err := db.c.Exec(`INSERT INTO users(username) VALUES (?);`, username)
		if err != nil {
			return id, false, err
		}

		res, err := result.LastInsertId()
		if err != nil {
			return id, false, err
		}

		return uint64(res), true, nil
	}
	if err != nil {
		return id, false, err
	}
	return id, false, nil
}

func (db *appdbimpl) SearchUsersByUsername(requestingUserId uint64, search string) (structs.UserList, error) {
	var result structs.UserList

	userSearchQuery := `
		SELECT * FROM users
		WHERE username LIKE '%' || ? || '%' AND
		NOT EXISTS (SELECT 1 FROM blocks WHERE 
			blockerUserId = users.id AND blockedUserId = ?)
		ORDER BY LENGTH(username) ASC
		LIMIT 24;
	`

	rows, err := db.c.Query(userSearchQuery, search, requestingUserId)

	if err != nil {
		return result, err
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
	return result, err
}

func (db *appdbimpl) GetUser(currentUserId uint64, requestedUserId uint64) (structs.User, error) {
	var result structs.User

	var (
		username          string
		numberOfFollowers uint64
		accountsFollowed  uint64
		numberOfPosts     uint64
		isBlocked         bool
		isFollowed        bool
	)

	userGetQuery := `
		SELECT
			u.username,
			COUNT(DISTINCT f1.followerUserId) as numberOfFollowers,
			COUNT(DISTINCT f2.followedUserId) as accountsFollowed,
			COUNT(DISTINCT p.postId) as numberOfPosts,
			MAX(CASE 
        		WHEN f3.followerUserId = ? THEN 1 
        		ELSE 0 
    		END) AS isFollowed,
			MAX(CASE 
        		WHEN b.blockerUserId = ? THEN 1 
        		ELSE 0 
    		END) AS isBlocked
		FROM users u
		LEFT JOIN follows f1 ON u.id = f1.followedUserId
		LEFT JOIN follows f2 ON u.id = f2.followerUserId
		LEFT JOIN follows f3 ON u.id = f3.followedUserId
		LEFT JOIN posts p ON u.id = p.userId
		LEFT JOIN blocks b ON ? = b.blockedUserId
		WHERE u.id = ?
	`

	err := db.c.QueryRow(userGetQuery, currentUserId, currentUserId, requestedUserId, requestedUserId).Scan(&username,
		&numberOfFollowers,
		&accountsFollowed,
		&numberOfPosts,
		&isFollowed,
		&isBlocked)

	if err != nil {
		return result, err
	}

	result = structs.User{
		Id:                requestedUserId,
		Username:          username,
		NumberOfFollowers: &numberOfFollowers,
		AccountsFollowed:  &accountsFollowed,
		NumberOfPosts:     &numberOfPosts,
		IsFollowed:        &isFollowed,
		IsBlocked:         &isBlocked,
	}

	return result, err
}

func (db *appdbimpl) CheckUsernameAvailability(newUsername string) (bool, error) {
	var response bool
	const checkUsernameAvailabilityQuery = "SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)"

	errors := db.c.QueryRow(checkUsernameAvailabilityQuery, newUsername).Scan(&response)

	return response, errors
}

func (db *appdbimpl) SetUsername(currentUserId uint64, newUsername string) error {
	const setUsernameQuery = "UPDATE users SET username = ? WHERE id = ?"

	_, err := db.c.Exec(setUsernameQuery, newUsername, currentUserId)

	return err
}
