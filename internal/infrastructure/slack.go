package infrastructure

import (
	"github.com/ikeisuke/slack-app-example/entity"
	"github.com/nlopes/slack"
)

type ISlack interface {
	PostMessage(channelID string, message *entity.SlackMessage) error
}

type Slack struct {
	upstream         *slack.Client
	DefaultChannelID string
}

func NewSlack(token string, channelID string) *Slack {
	return &Slack{
		upstream:         slack.New(token),
		DefaultChannelID: channelID,
	}
}

func (s *Slack) PostMessage(channelID string, message *entity.SlackMessage) error {
	if channelID == "" {
		channelID = s.DefaultChannelID
	}
	options := make([]slack.MsgOption, 0)
	options = append(options, slack.MsgOptionText(message.Text, false))
	if length := len(message.Attachments); length > 0 {
		attachments := make([]slack.Attachment, len(message.Attachments))
		for index, attachment := range message.Attachments {
			attachments[index] = slack.Attachment{
				Title: attachment.Title,
				Text:  attachment.Text,
			}
		}
		options = append(options, slack.MsgOptionAttachments(attachments...))
	}
	if message.ThreadTS != "" {
		options = append(options, slack.MsgOptionTS(message.ThreadTS))
	}
	_, _, err := s.upstream.PostMessage(channelID, options...)
	return err
}
