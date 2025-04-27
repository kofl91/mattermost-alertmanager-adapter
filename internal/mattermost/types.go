package mattermost

type MattermostMessage struct {
	Username    string       `json:"username,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Title   string   `json:"title,omitempty"`
	Text    string   `json:"text,omitempty"`
	Actions []Action `json:"actions,omitempty"`
	Fields  []Field  `json:"fields,omitempty"`
	Footer  string   `json:"footer,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Action struct {
	Id          string       `json:"id,omitempty"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	URL         string       `json:"url,omitempty"` // Keep for fallback
	Integration *Integration `json:"integration,omitempty"`
}

type Integration struct {
	URL     string                 `json:"url"`
	Context map[string]interface{} `json:"context,omitempty"`
}
