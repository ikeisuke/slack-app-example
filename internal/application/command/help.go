package command

import (
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/repository/command"
)

type HelpInput struct {
	SlashCommand string
}

type HelpInteraction struct {
	repository command.IHelpRepository
}

func NewHelpInteraction(repository command.IHelpRepository) *HelpInteraction {
	return &HelpInteraction{
		repository: repository,
	}
}

func (h *HelpInteraction) Run(input HelpInput) (*entity.SlackMessage, error) {
	message, err := h.repository.HelpMessage(command.HelpRepositoryInput{SlashCommand: input.SlashCommand})
	if err != nil {
		return nil, err
	}
	return &entity.SlackMessage{
		Blocks: []entity.SlackMessageBlock{
			{
				Type: "section",
				Text: &entity.SlackMessageTextObject{
					Type: "mrkdwn",
					Text: message,
				},
			},
		},
	}, nil
}
