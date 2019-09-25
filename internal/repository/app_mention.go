package repository

type AppMentionRepositoryInput struct {
	InnerEvent
	User    string `json:"user"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
	Channel string `json:"channel"`
	EventTS string `json:"event_ts"`
}

type AppMentionRepositoryOutput struct {
	ChannelID string
	Message   string
	Title     string
	Body      string
	TS        string
}

type IAppMentionRepository interface {
	Run(AppMentionRepositoryInput) (*AppMentionRepositoryOutput, error)
}

type AppMentionRepository struct {
}

func NewAppMentionRepository() *AppMentionRepository {
	return &AppMentionRepository{}
}

func (s *AppMentionRepository) Run(input AppMentionRepositoryInput) (*AppMentionRepositoryOutput, error) {
	message := &AppMentionRepositoryOutput{
		ChannelID: input.Channel,
		Message:   "app_mention reply!\nボットへのメンションのリプライ",
		Title:     "app_mention input text",
		//Body:      strings.Replace(input.Text, "@", "\\@", -1),
		Body:      input.Text,
		TS:        input.TS,
	}
	return message, nil
}
