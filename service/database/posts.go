package database

import (
	"time"

	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
)

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
	_, err := db.c.Exec(deletePostQuery, postId, currentUserId)

	return err
}

func (db *appdbimpl) CheckPostId(postId uint64) (bool, error) {
	var response bool
	const checkPostQuery = `SELECT EXISTS (SELECT 1 FROM posts WHERE postId = ?)`
	responseError := db.c.QueryRow(checkPostQuery, postId).Scan(&response)

	return response, responseError
}

func (db *appdbimpl) CheckLikePost(currentUserId uint64, postId uint64) (uint64, error) {
	var response uint64
	const checkLikeQuery = "SELECT COUNT(*) FROM likes WHERE likedPostId = ? AND userId = ?"
	responseError := db.c.QueryRow(checkLikeQuery, postId, currentUserId).Scan(&response)

	if responseError != nil {
		return 0, responseError
	}

	return response, responseError
}

func (db *appdbimpl) GetLikes(currentUserId uint64, postId uint64) (structs.UserList, error) {
	var result structs.UserList

	const getLikesQuery = "SELECT pLikes.id, pLikes.username FROM ( SELECT u.id, u.username FROM users u LEFT JOIN likes l ON l.userId = u.id WHERE l.likedPostId = ? ) pLikes LEFT JOIN blocks b ON b.blockerUserId = pLikes.id WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL"

	rows, errors := db.c.Query(getLikesQuery, postId, currentUserId)

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

func (db *appdbimpl) InsertLikePost(currentUserId uint64, postId uint64) error {
	const insertLikeQuery = "INSERT INTO likes (likedPostId, userId) VALUES (?, ?)"
	const deleteLikeQuery = "DELETE FROM likes WHERE likedPostId = ? AND userId = ?"
	const incrementPostLike = "UPDATE posts SET likes = likes + 1"
	_, err := db.c.Exec(insertLikeQuery, postId, currentUserId)

	if err == nil {
		_, incrementErr := db.c.Exec(incrementPostLike)
		if incrementErr != nil {
			_, deleteLikeError := db.c.Exec(deleteLikeQuery, postId, currentUserId)
			if deleteLikeError != nil {
				return deleteLikeError
			}
			return incrementErr
		}
	}

	return err
}

func (db *appdbimpl) DeleteLikePost(currentUserId uint64, postId uint64) error {
	const deleteLikeQuery = "DELETE FROM likes WHERE likedPostId = ? AND userId = ?"
	const decrementPostLike = "UPDATE posts SET likes = likes - 1"
	_, err := db.c.Exec(deleteLikeQuery, postId, currentUserId)

	if err == nil {
		_, incrementErr := db.c.Exec(decrementPostLike)
		if incrementErr != nil {
			return incrementErr
		}
	}

	return err
}

func (db *appdbimpl) GetCommentsPost(currentUserId uint64, postId uint64) (structs.CommentList, error) {
	var result structs.CommentList

	const getCommentsQuery = "SELECT pComments.commentId, pComments.postId, pComments.userId, u.username as authorUsername, pComments.text FROM ( SELECT c.commentId, c.postId, c.userId, c.text FROM comments c LEFT JOIN posts cPosts ON cPosts.postId = ? WHERE cPosts.postId = c.postId ) pComments LEFT JOIN blocks b ON b.blockerUserId = pComments.userId LEFT JOIN users u ON u.id = pComments.userId WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL"

	rows, errors := db.c.Query(getCommentsQuery, postId, currentUserId)

	if errors != nil {
		return result, errors
	}
	defer rows.Close()

	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.AuthorUsername, &comment.Text); err != nil {
			return result, err
		}
		result.Comments = append(result.Comments, comment)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, errors
}

func (db *appdbimpl) InsertCommentPost(currentUserId uint64, postId uint64, text string) error {
	const insertCommentQuery = "INSERT INTO comments (postId, userId, text) VALUES (?, ?, ?)"
	const deleteCommentQuery = "DELETE FROM comments WHERE commentId = ? AND userId = ?"
	const incrementPostLike = "UPDATE posts SET comments = comments + 1"
	result, err := db.c.Exec(insertCommentQuery, postId, currentUserId, text)

	if err == nil {
		commentId, errors := result.LastInsertId()

		if errors != nil {
			return errors
		}

		_, incrementErr := db.c.Exec(incrementPostLike)
		if incrementErr != nil {
			_, err := db.c.Exec(deleteCommentQuery, commentId, currentUserId)

			if err != nil {
				return err
			}

			return incrementErr
		}
	}

	return err
}

func (db *appdbimpl) DeleteCommentPost(currentUserId uint64, postId uint64, commentId uint64) error {
	const deleteCommentQuery = "DELETE FROM comments WHERE commentId = ? AND userId = ? AND postId = ?"
	const decrementPostLike = "UPDATE posts SET comments = comments - 1"
	_, err := db.c.Exec(deleteCommentQuery, commentId, currentUserId, postId)

	if err == nil {
		_, decrementErr := db.c.Exec(decrementPostLike)
		if decrementErr != nil {
			return decrementErr
		}
	}

	return err
}

func (db *appdbimpl) CheckCommentId(commentId uint64) (bool, error) {
	var response bool
	const checkCommentQuery = `SELECT EXISTS (SELECT 1 FROM comments WHERE commentId = ?)`
	responseError := db.c.QueryRow(checkCommentQuery, commentId).Scan(&response)

	return response, responseError
}

func (db *appdbimpl) GetPosts(currentUserId uint64, userToGetStream uint64, followedMode bool) (structs.Stream, error) {
	var result structs.Stream

	var getPostsQuery string

	if !followedMode {
		getPostsQuery = "SELECT p.*, CASE l.userId WHEN ? THEN TRUE ELSE FALSE END AS IsLiked FROM posts p LEFT JOIN likes l ON l.userId = ? AND l.likedPostId = p.postId  LEFT JOIN blocks b ON b.blockerUserId = p.userId WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL AND p.userId = ? ORDER BY p.uploadTime DESC"
	} else {
		getPostsQuery = "SELECT p.*, CASE l.userId WHEN ? THEN TRUE ELSE FALSE END AS IsLiked FROM posts p LEFT JOIN likes l ON l.userId = ? AND l.likedPostId = p.postId LEFT JOIN blocks b ON b.blockerUserId = p.userId LEFT JOIN follows f ON f.followerUserId = ? WHERE b.blockedUserId != ? OR b.blockedUserId IS NULL AND p.userId = f.followedUserId ORDER BY p.uploadTime DESC"
		userToGetStream = currentUserId
	}

	rows, errors := db.c.Query(getPostsQuery, currentUserId, currentUserId, currentUserId, userToGetStream)

	if errors != nil {
		return result, errors
	}
	defer rows.Close()

	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.Id, &post.Author, &post.Caption, &post.Likes, &post.Comments, &post.TimeOfCreation, &post.IsLiked); err != nil {
			return result, err
		}
		result.Stream = append(result.Stream, post)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, errors
}
