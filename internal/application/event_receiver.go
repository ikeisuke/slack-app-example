package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
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
	slack     infrastructure.ISlack
}

type EventReceiverOutput struct {
}

func NewEventReceiverInteraction(r repository.ISignatureRepository, e repository.IEventRepository, p presenter.IPresenter, s infrastructure.ISlack) *EventReceiverInteraction {
	return &EventReceiverInteraction{
		signature: r,
		event:     e,
		presenter: p,
		slack:     s,
	}
}

func (s *EventReceiverInteraction) Run(input *EventReceiverInput) string {
	var err error
	err = s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err == nil {
		var event repository.EventRepositoryInput
		if err = json.Unmarshal([]byte(input.Body), &event); err == nil {
			_, err = s.event.Run(event, func(eventType string, data []byte) error {
				switch eventType {
				case "app_mention":
					repo := repository.NewAppMentionRepository()
					app := NewEventAppMentionInteraction(repo, s.slack)
					return app.Run(data)
				default:
					return errors.New(fmt.Sprintf("unsupported inner event type: %s", eventType))
				}
				return nil
			})
		}
	}
	if err != nil {
		message := &infrastructure.Message{
			Text: err.Error(),
		}
		s.slack.PostMessage("", message)
	}
	return s.presenter.Output(input.Body, err)
}
