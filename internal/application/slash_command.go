package application

import (
	"errors"
	"github.com/ikeisuke/slack-app-example/internal/presenter"
	"github.com/ikeisuke/slack-app-example/internal/repository"
	"net/url"
	"strings"
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
	subCommand repository.ISubCommandRepository
	presenter  presenter.IPresenter
}

type SlackCommandOutput struct {
}

func NewSlashCommandInteraction(r repository.ISignatureRepository, s repository.ISubCommandRepository, p presenter.IPresenter) *SlackCommandInteraction {
	return &SlackCommandInteraction{
		signature:  r,
		subCommand: s,
		presenter:  p,
	}
}

func (s *SlackCommandInteraction) Run(input *SlashCommandInput) string {
	var body interface{}
	var err error
	err = s.signature.Verify(&repository.SignatureInput{
		Timestamp:        input.Timestamp,
		Signature:        input.Signature,
		SigningSecret:    input.SigningSecret,
		Body:             input.Body,
		SignatureVersion: input.SignatureVersion,
	})
	if err == nil {
		var parsed map[string]string
		parsed, err = parse(input.Body)
		if err == nil {
			input := repository.SubCommandRepositoryInput{
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
			}
			body, err = s.subCommand.Run(input)
		}
	}
	return s.presenter.Output(body, err)
}

func parse(qs string) (map[string]string, error) {
	input := make(map[string]string)
	for _, s := range strings.Split(qs, "&") {
		kv := strings.Split(s, "=")
		if len(kv) != 2 {
			return nil, errors.New("invalid request body detected")
		}
		key := kv[0]
		value, _ := url.QueryUnescape(kv[1])
		input[key] = value
	}
	return input, nil
}
