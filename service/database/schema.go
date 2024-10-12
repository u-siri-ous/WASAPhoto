package database

const (
	UsersSchema = ` 
		CREATE TABLE IF NOT EXISTS users (
			id              INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username        TEXT NOT NULL UNIQUE
		);
	`

	BlocksSchema = `
		CREATE TABLE IF NOT EXISTS blocks (
			blockedUserId INTEGER NOT NULL,
			blockerUserId INTEGER NOT NULL,
			PRIMARY KEY (blockedUserId, blockerUserId),
			FOREIGN KEY (blockerUserId) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (blockedUserId) REFERENCES users(id) ON DELETE CASCADE
		);
	`

	FollowsSchema = `
		CREATE TABLE IF NOT EXISTS follows (
			followerUserId INTEGER NOT NULL,
			followedUserId INTEGER NOT NULL,
			PRIMARY KEY (followerUserId, followedUserId),
			FOREIGN KEY (followerUserId) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (followedUserId) REFERENCES users(id) ON DELETE CASCADE
		);
	`

	PostsSchema = `
		CREATE TABLE IF NOT EXISTS posts (
			postId          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			userId        	INTEGER NOT NULL,
			caption			TEXT,
			likes			INTEGER,
			comments		INTEGER,
			uploadTime  	DATETIME NOT NULL,
		
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);
	`
)
