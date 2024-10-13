package requests

import "github.com/u-siri-ous/WASAPhoto/service/api/components/requests/utility"

type Comment struct {
	Text string `json:"text"`
}

func (request *Comment) IsValid() bool {
	return utility.CheckText(request.Text)
}
