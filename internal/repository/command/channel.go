package command

import "github.com/ikeisuke/slack-app-example/internal/infrastructure"

type IChannelRepository interface {
	List() (interface{}, error)
}

type ChannelRepository struct {
	infrastructure infrastructure.ISlack
}

func NewChannelRepository(infra infrastructure.ISlack) *ChannelRepository {
	return &ChannelRepository{
		infrastructure: infra,
	}
}

func (c *ChannelRepository) List() (interface{}, error) {
	return c.infrastructure.ChannelList()
}
