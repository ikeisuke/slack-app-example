package entity

import "encoding/json"

type EventWrapper struct {
	TeamID      string          `json:"team_id,omitempty"`
	APIAppID    string          `json:"api_app_id,omitempty"`
	Type        string          `json:"type,omitempty"`
	Challenge   string          `json:"challenge,omitempty"`
	AuthedUsers []string        `json:"authed_users,omitempty"`
	EventID     string          `json:"event_id,omitempty"`
	EventTime   int             `json:"event_time,omitempty"`
	Event       json.RawMessage `json:"event,omitempty"`
}

type EventBase struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype""`
}

type EventAppMention struct {
	EventBase
	User    string `json:"user"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
	Channel string `json:"channel"`
	EventTS string `json:"event_ts"`
}

type EventURLVerificationOutput struct {
	Challenge string `json:"challenge"`
}