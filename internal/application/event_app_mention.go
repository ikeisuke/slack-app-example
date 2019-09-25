package application

import (
	"encoding/json"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	"github.com/ikeisuke/slack-app-example/internal/repository"
)

type EventAppMentionInput struct {
}

type EventAppMentionInteraction struct {
	repository repository.IAppMentionRepository
	slack      infrastructure.ISlack
}

func NewEventAppMentionInteraction(a repository.IAppMentionRepository, s infrastructure.ISlack) *EventAppMentionInteraction {
	return &EventAppMentionInteraction{
		repository: a,
		slack:      s,
	}
}

func (e *EventAppMentionInteraction) Run(data []byte) error {
	var input repository.AppMentionRepositoryInput
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	if input.Subtype == "bot_message" {
		fmt.Println("skip bot message.")
		return nil
	}
	res, err := e.repository.Run(input)
	if err != nil {
		return err
	}
	message := &infrastructure.Message{
		Text: res.Message,
		Attachments: []infrastructure.MessageAttachment{
			{
				Title: res.Title,
				Text:  res.Body,
			},
		},
		ThreadTS: res.TS,
	}
	return e.slack.PostMessage(res.ChannelID, message)
}
