package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	application "github.com/ikeisuke/slack-app-example/internal/application/event"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	repository "github.com/ikeisuke/slack-app-example/internal/repository/event"
)

type EventRepositoryInput struct {
	OuterEvent entity.EventWrapper
}

type EventRepositoryOutput struct {
	Challenge string `json:"challenge"`
}

type IEventRepository interface {
	Run(input EventRepositoryInput) error
}
type EventRepository struct {
	infrastructure infrastructure.ISlack
}

func NewEventRepository(infra infrastructure.ISlack) *EventRepository {
	return &EventRepository{
		infrastructure: infra,
	}
}

func (s *EventRepository) Run(input EventRepositoryInput) error {
	rawInnerEvent := []byte(input.OuterEvent.Event)
	var inner entity.EventBase
	if err := json.Unmarshal(rawInnerEvent, &inner); err != nil {
		return err
	}
	switch inner.Type {
	case "app_mention":
		repo := repository.NewAppMentionRepository(s.infrastructure)
		app := application.NewAppMentionInteraction(repo)
		return app.Run(application.AppMentionInput{Data: rawInnerEvent})
	default:
		return errors.New(fmt.Sprintf("unsupported inner event type: %s", input.OuterEvent.Type))
	}
}
