package presenter

import (
	"encoding/json"
)

type ResponsePresenter struct {
}

func NewResponsePresenter() *ResponsePresenter {
	return &ResponsePresenter{}
}

type Response struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

func (r *ResponsePresenter) Output(data interface{}, err error) string {
	if err != nil {
		data = &Response{
			ResponseType: "ephemeral",
			Text:         "Sorry, that didn't work. Please try again. (" + err.Error() + ")",
		}
	}
	buf, _ := json.Marshal(data)
	return string(buf)
}
