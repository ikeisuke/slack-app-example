package presenter

type IPresenter interface {
	Output(interface{}) (string, error)
}
