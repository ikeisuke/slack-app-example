package presenter

import (
	"encoding/json"
)

type JSONPresenter struct {
}

func NewJSONPresenter() *JSONPresenter {
	return &JSONPresenter{}
}

func (r *JSONPresenter) Output(data interface{}) (string, error) {
	buf, err := json.Marshal(data)
	return string(buf), err
}
