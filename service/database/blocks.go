package database

import (
	"github.com/u-siri-ous/WASAPhoto/service/api/components/structs"
)

func (db *appdbimpl) BlockUser(currentUserId uint64, userToBlockId uint64) error {
	const blockUserQuery = "INSERT INTO blocks(blockedUserId, blockerUserId) VALUES (?, ?)"

	_, responseErrors := db.c.Exec(blockUserQuery, userToBlockId, currentUserId)

	return responseErrors
}

func (db *appdbimpl) UnblockUser(currentUserId uint64, userToUnblockId uint64) error {
	const unblockUserQuery = "DELETE FROM blocks WHERE blockedUserId = ? AND blockerUserId = ?"

	_, responseErrors := db.c.Exec(unblockUserQuery, userToUnblockId, currentUserId)

	return responseErrors
}

func (db *appdbimpl) CheckBlock(currentUserId uint64, userToCheck uint64) (bool, error) {
	var count int

	const checkBlockQuery = "SELECT COUNT(*) FROM blocks WHERE blockerUserId = ? AND blockedUserId = ?"

	responseErrors := db.c.QueryRow(checkBlockQuery, userToCheck, currentUserId).Scan(&count)

	if responseErrors != nil {
		return false, responseErrors
	}

	if count > 0 {
		return true, responseErrors
	} else {
		return false, responseErrors
	}
}

func (db *appdbimpl) ListOfBlocker(currentUserId uint64) (structs.UserList, error) {
	var result structs.UserList

	listOfBlockerQuery := `
		SELECT 
			u.id as userId,
			u.username
		FROM users u
		INNER JOIN blocks b
		ON u.id = b.blockedUserId
		WHERE b.blockerUserId = ?
	`

	response, responseErrors := db.c.Query(listOfBlockerQuery, currentUserId)

	if responseErrors != nil {
		return result, responseErrors
	}
	defer response.Close()

	for response.Next() {
		var user structs.User
		if responseErrors := response.Scan(&user.Id, &user.Username); responseErrors != nil {
			return result, responseErrors
		}
		result.Users = append(result.Users, user)
	}

	if responseErrors := response.Err(); responseErrors != nil {
		return result, responseErrors
	}
	return result, responseErrors
}
