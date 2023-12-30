package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookBody struct {
	Username  string          `json:"username"`
	AvatarURL string          `json:"avatar_url"`
	Content   string          `json:"content"`
	Embeds    []*WebhookEmbed `json:"embeds"`
}

type WebhookEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []WebhookEmbedField `json:"fields"`
	Footer      WebhookEmbedFooter  `json:"footer"`
	Thumbnail   URL                 `json:"thumbnail"`
	Image       URL                 `json:"image"`
}

type WebhookEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type WebhookEmbedFooter struct {
	Text string `json:"text"`
}

type URL struct {
	URL string `json:"url"`
}

type Webhook struct {
	URL  string
	Body WebhookBody
}

func (w *Webhook) Send() error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(w.Body); err != nil {
		return err
	}
	_, err := http.Post(w.URL, "application/json", &body)
	if err != nil {
		return err
	}
	return nil
}
