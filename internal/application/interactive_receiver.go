package application

import (
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
)

type InteractiveInput struct {
	Timestamp        int
	Signature        string
	SigningSecret    string
	Body             string
	SignatureVersion string
}

type InteractiveInteraction struct {
	signature  repository.ISignatureRepository
	repository repository.IInteractiveRepository
	presenter  presenter.IPresenter
}

func NewInteractivityInteraction(r repository.ISignatureRepository, i repository.IInteractiveRepository, p presenter.IPresenter) *InteractiveInteraction {
	return &InteractiveInteraction{
		signature:  r,
		repository: i,
		presenter:  p,
	}
}

func (s *InteractiveInteraction) Run(input InteractiveInput) error {
	err := s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err != nil {
		return err
	}
	data, err := parse(input.Body)
	if err != nil {
		return err
	}
	return s.repository.Run(repository.InteractiveRepositoryInput{Payload: data["payload"]})
}
