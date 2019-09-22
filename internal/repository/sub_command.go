package repository

import (
	"encoding/json"
)

type SubCommandRepositoryInput struct {
	ChannelID   string
	ChannelName string
	Command     string
	ResponseUrl string
	TeamDomain  string
	TeamID      string
	Text        string
	Token       string
	TriggerID   string
	UserID      string
	UserName    string
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
