package structs

type Comment struct {
	CommentId      uint64 `json:"commentId"`
	PostId         uint64 `json:"postId"`
	UserId         uint64 `json:"userId"`
	AuthorUsername string `json:"authorUsername"`
	Text           string `json:"text"`
}

type CommentList struct {
	Comments []Comment `json:"comments"`
}
