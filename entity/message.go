package entity


type SlackMessage struct {
	Text        string
	Attachments []SlackMessageAttachment
	ThreadTS    string
}

type SlackMessageAttachment struct {
	Title string
	Text  string
}