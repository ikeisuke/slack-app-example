package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
)

type EventReceiverInput struct {
	Timestamp        int
	Signature        string
	SigningSecret    string
	Body             string
	SignatureVersion string
}

type EventReceiverInteraction struct {
	signature  repository.ISignatureRepository
	repository repository.IEventRepository
	presenter  presenter.IPresenter
}

type EventReceiverOutput struct {
}

func NewEventReceiverInteraction(r repository.ISignatureRepository, e repository.IEventRepository, p presenter.IPresenter) *EventReceiverInteraction {
	return &EventReceiverInteraction{
		signature:  r,
		repository: e,
		presenter:  p,
	}
}

func (s *EventReceiverInteraction) Run(input *EventReceiverInput) (string, error) {
	verifyInput := repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	}
	if err := s.signature.Verify(verifyInput); err != nil {
		return "", err
	}
	var event entity.EventWrapper
	if err := json.Unmarshal([]byte(input.Body), &event); err == nil {
		return "", err
	}
	switch event.Type {
	case "url_verification":
		return s.presenter.Output(entity.EventURLVerificationOutput{
			Challenge: event.Challenge,
		})
	case "event_callback":
		return "", s.repository.Run(repository.EventRepositoryInput{OuterEvent: event})
	}
	return "", errors.New(fmt.Sprintf("unsupported outer event type: %s", event.Type))
}
