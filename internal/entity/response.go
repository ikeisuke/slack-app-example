package entity

type SlackAPIChannelListResponse = []SlackAPIChannelResponse

type SlackAPIChannelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
