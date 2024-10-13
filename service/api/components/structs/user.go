package structs

type User struct {
	Id                uint64  `json:"userId"`
	Username          string  `json:"username"`
	NumberOfFollowers *uint64 `json:"numberOfFollowers"`
	AccountsFollowed  *uint64 `json:"accountsFollowed"`
	NumberOfPosts     *uint64 `json:"numberOfPosts"`
	IsBlocked         *bool   `json:"isBlocked"`
	IsFollowed        *bool   `json:"isFollowed"`
}

type UserReduced struct {
	Id       uint64 `json:"userId"`
	Username string `json:"username"`
}

type UserList struct {
	Users []UserReduced `json:"users"`
}
