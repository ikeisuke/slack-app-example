package event

import (
	"errors"
	"fmt"
	"github.com/ikeisuke/slack-app-example/internal/event/application"
	"github.com/ikeisuke/slack-app-example/internal/event/presenter"
	"github.com/ikeisuke/slack-app-example/internal/event/repository"
)

type EventWrapper struct {
	Token       string      `json:"token,omitempty"`
	TeamID      string      `json:"team_id,omitempty"`
	APIAppID    string      `json:"api_app_id,omitempty"`
	Event       interface{} `json:"event,omitempty"`
	Type        string      `json:"type,omitempty"`
	Challenge   string      `json:"challenge,omitempty"`
	AuthedUsers []string    `json:"authed_users,omitempty"`
	EventID     string      `json:"event_id,omitempty"`
	EventTime   int         `json:"event_time,omitempty"`
}

func Router(event EventWrapper) (interface{}, error) {
	if event.Type == "url_verification" {
		repo := repository.NewURLVerificationRepository()
		pres := &presenter.URLVerificationPresenter{}
		app := application.NewURLVerificationInteraction(repo, pres)
		return app.Run(event.Challenge)
	}
	return nil, errors.New(fmt.Sprintf("unknown event type: %s", event.Type))
}
