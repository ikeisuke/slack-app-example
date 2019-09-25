package application

import (
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"strings"
)

type EventReceiverInput struct {
	Timestamp        int
	Signature        string
	SigningSecret    string
	Body             string
	SignatureVersion string
}

type EventReceiverInteraction struct {
	signature      repository.ISignatureRepository
	event          repository.IEventRepository
	presenter      presenter.IPresenter
	infrastructure infrastructure.ISlack
}

type EventReceiverOutput struct {
}

func NewEventReceiverInteraction(r repository.ISignatureRepository, e repository.IEventRepository, p presenter.IPresenter, i infrastructure.ISlack) *EventReceiverInteraction {
	return &EventReceiverInteraction{
		signature:      r,
		event:          e,
		presenter:      p,
		infrastructure: i,
	}
}

func (s *EventReceiverInteraction) Run(input *EventReceiverInput) string {
	var body interface{}
	var err error
	err = s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err == nil {
		body, err = s.event.Run(input.Body)
	}
	if err != nil {
		message := &infrastructure.Message{
			Text: err.Error(),
			Attachments: []infrastructure.MessageAttachment{
				{
					Title: "request body",
					Text:  strings.Replace(input.Body, "@", "\\@", -1),
				},
			},
		}
		s.infrastructure.PostMessage("", message)
	}
	return s.presenter.Output(body, err)
}
