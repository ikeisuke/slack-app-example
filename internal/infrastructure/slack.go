package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
)

type ISlack interface {
	PostMessage(channelID string, message *entity.SlackMessage) error
	ChannelList() (interface{}, error)
	PostToResponseURL(url string, data interface{}) error
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

func (s *Slack) ChannelList() (interface{}, error) {
	return s.upstream.GetChannels(false)
}

type PostToResponseURLResponse struct {
	OK    string `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (s *Slack) PostToResponseURL(url string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return err
	}
	parsed := &PostToResponseURLResponse{}
	if err := json.Unmarshal(body, res); err != nil {
		return err
	}
	if parsed.OK != "ok" {
		return errors.New(parsed.Error)
	}
	return nil
}
