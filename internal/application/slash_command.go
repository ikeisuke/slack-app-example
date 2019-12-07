package application

import (
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"net/url"
)

type SlashCommandInput struct {
	Timestamp        int
	Signature        string
	SigningSecret    string
	Body             string
	SignatureVersion string
}

type SlackCommandInteraction struct {
	signature  repository.ISignatureRepository
	subCommand repository.ICommandRepository
	presenter  presenter.IPresenter
}

type SlackCommandOutput struct {
}

func NewSlashCommandInteraction(r repository.ISignatureRepository, s repository.ICommandRepository, p presenter.IPresenter) *SlackCommandInteraction {
	return &SlackCommandInteraction{
		signature:  r,
		subCommand: s,
		presenter:  p,
	}
}

func (s *SlackCommandInteraction) Run(input *SlashCommandInput) (string, error) {
	err := s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err != nil {
		return s.presenter.Output(errorData(err))
	}
	parsed, err := parse(input.Body)
	if err != nil {
		return s.presenter.Output(errorData(err))
	}
	body, err := s.subCommand.Run(repository.CommandRepositoryInput{
		ChannelID:   parsed["channel_id"],
		ChannelName: parsed["channel_name"],
		Command:     parsed["command"],
		ResponseUrl: parsed["response_url"],
		TeamDomain:  parsed["team_domain"],
		TeamID:      parsed["team_id"],
		Text:        parsed["text"],
		Token:       parsed["token"],
		TriggerID:   parsed["trigger_id"],
		UserID:      parsed["user_id"],
		UserName:    parsed["user_name"],
	})
	if err != nil {
		return s.presenter.Output(errorData(err))
	}
	return s.presenter.Output(body)
}

func errorData(err error) entity.SlackMessage {
	return entity.SlackMessage{
		ResponseType: "ephemeral",
		Text:         "Sorry, that didn't work. Please try again. (" + err.Error() + ")",
	}
}

func parse(qs string) (map[string]string, error) {
	result := make(map[string]string)
	parsed, err := url.ParseQuery(qs)
	if err != nil {
		return nil, err
	}
	for key, value := range parsed {
		result[key] = value[0]
	}
	return result, nil
}
