package repository

import (
	"encoding/json"
	"github.com/ikeisuke/slack-app-example/internal/event"
)

type IEventRepository interface {
	Run(input string) (interface{}, error)
}
type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (s *EventRepository) Run(input string) (interface{}, error) {
	e := event.EventWrapper{}
	err := json.Unmarshal([]byte(input), &e)
	if err != nil {
		return nil, err
	}
	data, err := event.Router(e)
	if err != nil {
		return nil, err
	}
	return data, nil
}
