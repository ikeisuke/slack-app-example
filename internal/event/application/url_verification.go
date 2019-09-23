package application

import (
	"github.com/ikeisuke/slack-app-example/internal/event/presenter"
	"github.com/ikeisuke/slack-app-example/internal/event/repository"
)

type IURLVerificationInteraction interface {
	Run(s string) (interface{}, error)
}

type URLVerificationInteraction struct {
	repository repository.IURLVerificationRepository
	presenter  presenter.IURLVerificationPresenter
}

func NewURLVerificationInteraction(r repository.IURLVerificationRepository, p presenter.IURLVerificationPresenter) *URLVerificationInteraction {
	return &URLVerificationInteraction{
		repository: r,
		presenter:  p,
	}
}

func (u *URLVerificationInteraction) Run(s string) (interface{}, error) {
	return u.presenter.Format(u.repository.Challenge(s)), nil
}
