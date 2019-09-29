package repository

import (
	"encoding/json"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
)

type InteractiveRepositoryInput struct {
	Payload string
}

type IInteractiveRepository interface {
	Run(input InteractiveRepositoryInput) error
}
type InteractiveRepository struct {
	infrastructure infrastructure.ISlack
}

func NewInteractiveRepository(infra infrastructure.ISlack) *InteractiveRepository {
	return &InteractiveRepository{
		infrastructure: infra,
	}
}

func (s *InteractiveRepository) Run(input InteractiveRepositoryInput) error {
	var payload entity.Interactivity
	if err := json.Unmarshal([]byte(input.Payload), &payload); err != nil {
		return err
	}
	message := &entity.SlackMessage{
		ReplaceOriginal: true,
		Blocks: []entity.SlackMessageBlock{
			{
				Type: "section",
				Text: &entity.SlackMessageTextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("<@%s> push %s button", payload.User.ID, payload.Actions[0].Text.Text),
				},
			},
		},
	}
	return s.infrastructure.PostToResponseURL(payload.ResponseURL, message)
}
