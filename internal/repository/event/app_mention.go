package event

import (
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
)

type IAppMentionRepository interface {
	Run(input entity.EventAppMention) error
}

type AppMentionRepository struct {
	infrastructure infrastructure.ISlack
}

func NewAppMentionRepository(infra infrastructure.ISlack) *AppMentionRepository {
	return &AppMentionRepository{
		infrastructure: infra,
	}
}

func (s *AppMentionRepository) Run(input entity.EventAppMention) error {
	return s.infrastructure.PostMessage(input.Channel, &entity.SlackMessage{
		Text: "app_mention reply!\nボットへのメンションのリプライ",
		Attachments: []entity.SlackMessageAttachment{
			{
				Title: "app_mention input text",
				Text:  input.Text,
			},
		},
		ThreadTS: input.TS,
	})
}
