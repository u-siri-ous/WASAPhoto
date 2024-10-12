/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)

	//Authorization
	CheckUserId(userId uint64) (bool, error)

	//Login
	InsertUser(username string) (uint64, bool, error)

	//Search User
	SearchUsersByUsername(requesterUserId uint64, search string) (structs.UserList, error)
	GetUser(currentUserId uint64, requestedUserId uint64) (structs.User, error)

	//Ban
	BlockUser(userId uint64, userToBlockId uint64) error
	UnblockUser(userId uint64, userToBlockId uint64) error
	CheckBlock(userId uint64, userToCheckId uint64) (bool, error)
	ListOfBlocker(userid uint64) (structs.UserList, error)

	//Follow
	FollowUser(userId uint64, userToFollowId uint64) error
	UnfollowUser(userId uint64, userToUnfollowId uint64) error

	//Posts
	CreatePost(userId uint64, caption string, uploadTime time.Time) (uint64, error)
	DeletePost(currentUserId uint64, postId uint64) error
	CheckPostId(postId uint64) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	TableMapping := map[string]string{
		"users":   UsersSchema,
		"blocks":  BlocksSchema,
		"follows": FollowsSchema,
		"posts":   PostsSchema,
	}

	for tableName, sqlStmt := range TableMapping {
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name= ? ;`, tableName).Scan(&tableName)

		if errors.Is(err, sql.ErrNoRows) {
			_, err = db.Exec(sqlStmt)

			if err != nil {
				return nil, fmt.Errorf("error creating database structure.\n%s -> %w", tableName, err)
			}
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
