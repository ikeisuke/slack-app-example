package command

import (
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	repository "github.com/ikeisuke/slack-app-example/internal/repository/command"
)

type ChannelInteraction struct {
	repository repository.IChannelRepository
}

type ChannelListInput struct {
}

func NewChannelInteraction(repo repository.IChannelRepository) *ChannelInteraction {
	return &ChannelInteraction{
		repository: repo,
	}
}

func (c *ChannelInteraction) List(input ChannelListInput) (*entity.SlackMessage, error) {
	channels, err := c.repository.List()
	if err != nil {
		return nil, err
	}
	elements := make([]entity.SlackMessageBlockActionElement, len(channels), len(channels))
	for i, channel := range channels {
		elements[i] = entity.SlackMessageBlockActionElement{
			Type:     "button",
			ActionID: fmt.Sprintf("channel_detail_%s", channel.ID),
			Text: &entity.SlackMessageTextObject{
				Type: "plain_text",
				Text: channel.Name,
			},
		}
	}
	return &entity.SlackMessage{
		ResponseType: "in_channel",
		Blocks: []entity.SlackMessageBlock{
			{
				Type:     "actions",
				Elements: elements,
			},
		},
	}, nil
}
