package command

import (
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
)

type IChannelRepository interface {
	List() (entity.SlackAPIChannelListResponse, error)
}

type ChannelRepository struct {
	infrastructure infrastructure.ISlack
}

func NewChannelRepository(infra infrastructure.ISlack) *ChannelRepository {
	return &ChannelRepository{
		infrastructure: infra,
	}
}

func (c *ChannelRepository) List() (entity.SlackAPIChannelListResponse, error) {
	return c.infrastructure.ChannelList()
}
