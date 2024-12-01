package structs

type Comment struct {
	CommentId uint64 `json:"commentId"`
	PostId    uint64 `json:"postId"`
	UserId    uint64 `json:"userId"`
	Text      string `json:"text"`
}

type CommentList struct {
	Comments []Comment `json:"comments"`
}
