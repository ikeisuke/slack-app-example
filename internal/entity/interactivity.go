package entity

type Interactivity struct {
	Type        string                 `json:"type"`
	Team        InteractivityTeam      `json:"team"`
	User        InteractivityUser      `json:"user"`
	APIAppID    string                 `json:"api_app_id"`
	Container   InteractivityContainer `json:"container"`
	TriggerID   string                 `json:"trigger_id"`
	Channel     InteractivityChannel   `json:"channel"`
	ResponseURL string                 `json:"response_url"`
	Actions     []InteractivityAction  `json:"actions"`
}

type InteractivityTeam struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
}

type InteractivityUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	TeamID   string `json:"team_id"`
}

type InteractivityContainer struct {
	Type        string `json:"type"`
	MessageTS   string `json:"message_ts"`
	ChannelID   string `json:"channel_id"`
	IsEphemeral bool   `json:"is_ephemeral"`
}

type InteractivityChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type InteractivityAction struct {
	ActionID string                  `json:"action_id"`
	BlockID  string                  `json:"block_id"`
	Text     InteractivityActionText `json:"text"`
	Type     string                  `json:"type"`
	ActionTS string                  `json:"action_ts"`
}

type InteractivityActionText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	emoji bool   `json:"emoji"`
}
