package requests

import "github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"

type Username struct {
	Username string `json:"username"`
}

func (request *Username) IsValid() bool {
	return utility.CheckUsername(request.Username)
}
