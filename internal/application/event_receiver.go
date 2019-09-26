package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ikeisuke/slack-app-example/entity"
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
	signature repository.ISignatureRepository
	event     repository.IEventRepository
	presenter presenter.IPresenter
}

type EventReceiverOutput struct {
}

func NewEventReceiverInteraction(r repository.ISignatureRepository, e repository.IEventRepository, p presenter.IPresenter) *EventReceiverInteraction {
	return &EventReceiverInteraction{
		signature: r,
		event:     e,
		presenter: p,
	}
}

func (s *EventReceiverInteraction) Run(input *EventReceiverInput) string {
	var err error
	var data interface{}
	err = s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err == nil {
		var event entity.EventWrapper
		if err = json.Unmarshal([]byte(input.Body), &event); err == nil {
			switch event.Type {
			case "url_verification":
				data = entity.EventURLVerificationOutput{
					Challenge: event.Challenge,
				}
			case "event_callback":
				err = s.event.Run(repository.EventRepositoryInput{OuterEvent: event})
			default:
				err = errors.New(fmt.Sprintf("unsupported outer event type: %s", event.Type))
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return s.presenter.Output(data, err)
}
