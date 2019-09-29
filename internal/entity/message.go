package entity

type SlackMessage struct {
	ResponseType    string                   `json:"response_type,omitempty"`
	Text            string                   `json:"text,omitempty"`
	Attachments     []SlackMessageAttachment `json:"attachments,omitempty"`
	Blocks          []SlackMessageBlock      `json:"blocks,omitempty"`
	ThreadTS        string                   `json:"thread_ts,omitempty"`
	ReplaceOriginal bool                     `json:"replace_original,omitempty"`
}

type SlackMessageAttachment struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

type SlackMessageBlock struct {
	Type     string                           `json:"type,omitempty"`
	Text     *SlackMessageTextObject          `json:"text,omitempty"`
	Elements []SlackMessageBlockActionElement `json:"elements,omitempty"`
}

type SlackMessageTextObject struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type SlackMessageBlockActionElement struct {
	Type     string                  `json:"type,omitempty"`
	ActionID string                  `json:"action_id,omitempty"`
	Text     *SlackMessageTextObject `json:"text,omitempty"`
	URL      string                  `json:"url,omitempty"`
}
