package repository

import (
	"errors"
	"fmt"
	application "github.com/ikeisuke/slack-app-example/internal/application/command"
	"github.com/ikeisuke/slack-app-example/internal/entity"
	"github.com/ikeisuke/slack-app-example/internal/infrastructure"
	repository "github.com/ikeisuke/slack-app-example/internal/repository/command"
	"strings"
)

type CommandRepositoryInput struct {
	ChannelID   string `json:"channel_id,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	Command     string `json:"command,omitempty"`
	ResponseUrl string `json:"response_url,omitempty"`
	TeamDomain  string `json:"team_domain,omitempty"`
	TeamID      string `json:"team_id,omitempty"`
	Text        string `json:"text,omitempty"`
	Token       string `json:"token,omitempty"`
	TriggerID   string `json:"trigger_id,omitempty"`
	UserID      string `json:"user_id,omitempty"`
	UserName    string `json:"user_name,omitempty"`
}

type ICommandRepository interface {
	Run(input CommandRepositoryInput) (*entity.SlackMessage, error)
}
type CommandRepository struct {
	infrastructure infrastructure.ISlack
}

func NewCommandRepository(infra infrastructure.ISlack) *CommandRepository {
	return &CommandRepository{
		infrastructure: infra,
	}
}

func (s *CommandRepository) Run(input CommandRepositoryInput) (*entity.SlackMessage, error) {
	command := ""
	subcommand := ""
	args := parseArgs(input.Text)
	if len(args) > 0 {
		command = strings.ToLower(args[0])
		args = args[1:]
	}
	if len(args) > 0 {
		subcommand = strings.ToLower(args[0])
		args = args[1:]
	}
	switch command {
	case "help", "h", "":
		repo := repository.NewHelpRepository()
		app := application.NewHelpInteraction(repo)
		return app.Run(application.HelpInput{SlashCommand: input.Command})
	case "channel", "c":
		repo := repository.NewChannelRepository(s.infrastructure)
		app := application.NewChannelInteraction(repo)
		if subcommand == "list" {
			return app.List(application.ChannelListInput{})
		}
	}
	return nil, errors.New(fmt.Sprintf("Unsuppoerted command: %s %s %s", input.Command, command, subcommand))
}

func parseArgs(s string) []string {
	buf := []byte(strings.Trim(s, " "))
	result := make([]string, 0, 0)
	quotedMarks := []byte("\"'")
	separator := []byte(" ")

	var quoting bool
	var currentQuotedMarker byte
	index := 0
	for i, b := range buf {
		if !quoting && find(b, quotedMarks) {
			currentQuotedMarker = b
			quoting = true
			index = i + 1
			continue
		}
		if quoting && currentQuotedMarker == b {
			currentQuotedMarker = 0
			quoting = false
			result = append(result, string(buf[index:i]))
			index = i + 1
			continue
		}
		if !quoting && find(b, separator) {
			if i == index {
				index = i + 1
				continue
			}
			result = append(result, string(buf[index:i]))
			index = i + 1
			continue
		}
	}
	last := buf[index:]
	if len(last) > 0 {
		result = append(result, string(last))
	}
	return result
}

func find(b byte, list []byte) bool {
	for _, v := range list {
		if v == b {
			return true
		}
	}
	return false
}
