package database

import "time"

func (db *appdbimpl) CreatePost(currentUserId uint64, caption string, uploadTime time.Time) (uint64, error) {
	const createPostQuery = "INSERT INTO posts (userId, caption, likes, comments, uploadTime) VALUES (?, ?, 0, 0, ?)"

	result, err := db.c.Exec(createPostQuery, currentUserId, caption, uploadTime)
	if err != nil {
		return 0, err
	}

	postId, parseIdError := result.LastInsertId()
	if parseIdError != nil {
		return 0, parseIdError
	}

	return uint64(postId), nil
}

func (db *appdbimpl) DeletePost(currentUserId uint64, postId uint64) error {
	const deletePostQuery = "DELETE FROM posts WHERE postId = ? AND userId = ?"
	_, err := db.c.Exec(deletePostQuery, currentUserId, postId)

	return err
}

func (db *appdbimpl) CheckPostId(postId uint64) (bool, error) {
	var response bool
	const checkPostQuery = `SELECT TRUE FROM posts WHERE postId = ?`
	responseError := db.c.QueryRow(checkPostQuery, postId).Scan(&response)

	return response, responseError
}
