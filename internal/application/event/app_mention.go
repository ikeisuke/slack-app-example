package event

import (
	"encoding/json"
	"errors"
	"github.com/ikeisuke/slack-app-example/entity"
	repository "github.com/ikeisuke/slack-app-example/internal/repository/event"
)

type AppMentionInput struct {
	Data []byte
}

type AppMentionInteraction struct {
	repository repository.IAppMentionRepository
}

func NewAppMentionInteraction(a repository.IAppMentionRepository) *AppMentionInteraction {
	return &AppMentionInteraction{
		repository: a,
	}
}

func (e *AppMentionInteraction) Run(input AppMentionInput) error {
	var event entity.EventAppMention
	if err := json.Unmarshal(input.Data, &event); err != nil {
		return err
	}
	if event.Subtype == "bot_message" {
		return errors.New("skip bot message")
	}
	return e.repository.Run(event)
}
