package structs

import "time"

type Post struct {
	Id             uint64    `json:"photoId"`
	Author         string    `json:"author"`
	AuthorId       uint64    `json:"authorId"`
	Caption        string    `json:"caption"`
	Likes          uint64    `json:"numberOfLikes"`
	Comments       uint64    `json:"numberOfComments"`
	TimeOfCreation time.Time `json:"timeOfCreation"`
	IsLiked        bool      `json:"isliked"`
}

type Stream struct {
	Stream []Post `json:"stream"`
}
