package presenter

type IURLVerificationPresenter interface {
	Format(s string) map[string]string
}

type URLVerificationPresenter struct {
}

func (u *URLVerificationPresenter) Format(s string) map[string]string {
	return map[string]string{
		"challenge": s,
	}
}
