package repository

import (
	"encoding/json"
)

type SubCommandRepositoryInput struct {
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

type ISubCommandRepository interface {
	Run(input SubCommandRepositoryInput) (interface{}, error)
}
type SubCommandRepository struct{}

func NewSubCommandRepository() *SubCommandRepository {
	return &SubCommandRepository{}
}

func (s *SubCommandRepository) Run(input SubCommandRepositoryInput) (interface{}, error) {
	data, _ := json.MarshalIndent(input, "", "  ")
	return map[string]string{
		"text": string(data),
	}, nil
}
