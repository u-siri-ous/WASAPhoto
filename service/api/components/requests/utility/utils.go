package utility

import (
	"regexp"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*[\w]$`)

func CheckUsername(username string) bool {
	return 3 <= len(username) && len(username) <= 25 && usernameRegex.MatchString(username)
}

func CheckText(text string) bool {
	return 1 <= len(text) && len(text) <= 255
}
