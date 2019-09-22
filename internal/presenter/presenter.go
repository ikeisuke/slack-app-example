package presenter

type IPresenter interface {
	Output(interface{}, error) string
}
