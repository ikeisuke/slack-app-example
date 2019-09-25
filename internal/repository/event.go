package repository

import (
	"encoding/json"
	"errors"
	"fmt"
)

type EventRepositoryInput struct {
	TeamID      string          `json:"team_id,omitempty"`
	APIAppID    string          `json:"api_app_id,omitempty"`
	Type        string          `json:"type,omitempty"`
	Challenge   string          `json:"challenge,omitempty"`
	AuthedUsers []string        `json:"authed_users,omitempty"`
	EventID     string          `json:"event_id,omitempty"`
	EventTime   int             `json:"event_time,omitempty"`
	Event       json.RawMessage `json:"event,omitempty"`
}

type InnerEvent struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype""`
}

type EventRepositoryOutput struct {
	Challenge string `json:"challenge"`
}

type IEventRepository interface {
	Run(input EventRepositoryInput, callback func(string, []byte) error) (interface{}, error)
}
type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (s *EventRepository) Run(input EventRepositoryInput, callback func(string, []byte) error) (interface{}, error) {
	if input.Type == "url_verification" {
		return &EventRepositoryOutput{
			Challenge: input.Challenge,
		}, nil
	} else if input.Type == "event_callback" {
		event := []byte(input.Event)
		var inner InnerEvent
		if err := json.Unmarshal(event, &inner); err != nil {
			return nil, err
		}
		if err := callback(inner.Type, event); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, errors.New(fmt.Sprintf("unknown event type: %s", input.Type))
}
