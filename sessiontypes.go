package wabaapi

import (
	"encoding/json"
	"net/url"
)

type OutboundMessage struct {
	Channel        string
	Destination    string
	Source         string
	SourceName     string
	DisablePreview bool
	Validate       bool
}

func (om *OutboundMessage) defaultValues() url.Values {
	values := url.Values{}
	values.Add("channel", om.Channel)
	values.Add("destination", om.Destination)
	values.Add("source", om.Source)
	values.Add("src.name", om.SourceName)
	values.Add("disablePreview", "true")
	return values
}

func (om *OutboundMessage) Text(text string) url.Values {
	values := om.defaultValues()
	msg := map[string]string{
		"type": "text",
		"text": text,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values
}

func (om *OutboundMessage) Image(originalURL string, previewURL string) url.Values {
	values := om.defaultValues()
	msg := map[string]string{
		"type":        "image",
		"originalUrl": originalURL,
		"previewUrl":  previewURL,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values
}

func (om *OutboundMessage) File(url string, filename string) url.Values {
	values := om.defaultValues()
	msg := map[string]string{
		"type":     "file",
		"url":      url,
		"filename": filename,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values
}

func (om *OutboundMessage) Audio(url string) url.Values {
	values := om.defaultValues()
	msg := map[string]string{
		"type": "audio",
		"url":  url,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values
}

func (om *OutboundMessage) Video(url string, caption string) url.Values {
	values := om.defaultValues()
	msg := map[string]string{
		"type":    "video",
		"url":     url,
		"caption": caption,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values
}

func (om *OutboundMessage) ListMessage(lm ListMessage) url.Values {
	values := om.defaultValues()
	txt, _ := json.Marshal(lm)
	values.Add("message", string(txt))
	return values

}

type ListMessage struct {
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	MsgID        string     `json:"msgid,omitempty"`
	GlobalButton string     `json:"-"`
	Items        []ListItem `json:"items"`
}

func (lm ListMessage) MarshalJSON() ([]byte, error) {
	type Alias ListMessage
	type tmplm struct {
		Alias
		Type   string              `json:"type"`
		Button []map[string]string `json:"globalButtons"`
	}

	btn := []map[string]string{
		{
			"type":  "text",
			"title": lm.GlobalButton,
		},
	}

	tmp := tmplm{
		Alias:  Alias(lm),
		Type:   "list",
		Button: btn,
	}

	return json.Marshal(tmp)
}

type ListItem struct {
	Title   string           `json:"title"`
	Options []ListItemOption `json:"options"`
}

type ListItemOption struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	PostbackText string `json:"postbackText"`
}

func (li ListItemOption) MarshalJSON() ([]byte, error) {
	type Alias ListItemOption
	type tmpli struct {
		Alias
		Type string `json:"type"`
	}

	tmp := tmpli{
		Alias: Alias(li),
		Type:  "text",
	}

	return json.Marshal(tmp)
}
