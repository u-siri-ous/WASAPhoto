package structs

import "time"

type Post struct {
	Id             uint64    `json:"photoId"`
	Author         uint64    `json:"author"`
	Caption        string    `json:"caption"`
	Likes          uint64    `json:"numberOfLikes"`
	Comments       uint64    `json:"numberOfComments"`
	TimeOfCreation time.Time `json:"timeOfCreation"`
	IsLiked        bool      `json:"isliked"`
}
