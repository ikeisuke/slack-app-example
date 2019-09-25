package presenter

import "encoding/json"

type SimplePresenter struct {
}

func NewSimplePresenter() *SimplePresenter {
	return &SimplePresenter{}
}

func (t *SimplePresenter) Output(data interface{}, err error) string {
	if err != nil {
		data = map[string]string{
			"err": err.Error(),
		}
	}
	buf, _ := json.Marshal(data)
	return string(buf)
}

