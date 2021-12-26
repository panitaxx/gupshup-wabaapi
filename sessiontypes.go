package wabaapi

import (
	"encoding/json"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

//OutboundMessage is the basic structure for creating reply messages
//Call this structure with the appropriate method to create a reply message
//Limited validation is performed on the structure
type OutboundMessage struct {
	Channel        string
	Destination    string
	Source         string
	SourceName     string
	DisablePreview bool
	DoNotValidate  bool
}

func (om *OutboundMessage) Validate() error {
	if om.DoNotValidate {
		return nil
	}

	return validation.ValidateStruct(&om,
		validation.Field(&om.Channel, validation.Required),
		validation.Field(&om.Destination, validation.Required, is.E164),
		validation.Field(&om.Source, validation.Required),
		validation.Field(&om.SourceName, validation.Required),
	)
}

func (om *OutboundMessage) defaultValues() (url.Values, error) {
	if err := om.Validate(); err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Add("channel", om.Channel)
	values.Add("destination", om.Destination)
	values.Add("source", om.Source)
	values.Add("src.name", om.SourceName)
	values.Add("disablePreview", "true")
	return values, nil
}

//Text creates a text message
func (om *OutboundMessage) Text(text string) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	msg := map[string]string{
		"type": "text",
		"text": text,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values, nil
}

//Image creates an image message
func (om *OutboundMessage) Image(originalURL string, previewURL string) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	msg := map[string]string{
		"type":        "image",
		"originalUrl": originalURL,
		"previewUrl":  previewURL,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values, nil
}

//File creates a file message
func (om *OutboundMessage) File(url string, filename string) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	msg := map[string]string{
		"type":     "file",
		"url":      url,
		"filename": filename,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values, nil
}

//Audio creates an audio message
func (om *OutboundMessage) Audio(url string) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	msg := map[string]string{
		"type": "audio",
		"url":  url,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values, nil
}

//Video creates a video message
func (om *OutboundMessage) Video(url string, caption string) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	msg := map[string]string{
		"type":    "video",
		"url":     url,
		"caption": caption,
	}
	txt, _ := json.Marshal(msg)
	values.Add("message", string(txt))
	return values, nil
}

//Creates an interactive list message
func (om *OutboundMessage) ListMessage(lm ListMessage) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	txt, _ := json.Marshal(lm)
	values.Add("message", string(txt))
	return values, nil

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

func (om *OutboundMessage) QuickReplyText(text QuickReplyText) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	txt, _ := json.Marshal(text)
	values.Add("message", string(txt))
	return values, nil
}

func (om *OutboundMessage) QuickReplyImage(qri QuickReplyImage) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	txt, _ := json.Marshal(qri)
	values.Add("message", string(txt))
	return values, nil
}

func (om *OutboundMessage) QuickReplyDocument(qrd QuickReplyDocument) (url.Values, error) {
	values, err := om.defaultValues()
	if err != nil {
		return nil, err
	}
	txt, _ := json.Marshal(qrd)
	values.Add("message", string(txt))
	return values, nil
}

type QuickReplyOption string

func (qro *QuickReplyOption) MarshalJSON() ([]byte, error) {
	type tmpli struct {
		Opt  string `json:"text"`
		Type string `json:"type"`
	}

	tmp := tmpli{
		Opt:  string(*qro),
		Type: "text",
	}

	return json.Marshal(tmp)
}

type QuickReplyImage struct {
	MsgID   string
	URL     string
	Text    string
	Caption string
	Options []QuickReplyOption
}

func (qri *QuickReplyImage) MarshalJSON() ([]byte, error) {

	type TContent struct {
		Type    string `json:"type"`
		URL     string `json:"url"`
		Text    string `json:"text"`
		Caption string `json:"caption"`
	}

	tmpli := struct {
		MsgID   string             `json:"msgid"`
		Content TContent           `json:"content"`
		Options []QuickReplyOption `json:"options"`
	}{
		MsgID:   qri.MsgID,
		Options: qri.Options,
		Content: TContent{
			Type:    "image",
			URL:     qri.URL,
			Text:    qri.Text,
			Caption: qri.Caption,
		},
	}

	return json.Marshal(tmpli)
}

type QuickReplyText struct {
	MsgID   string `json:"msgid"`
	Header  string `json:"header"`
	Text    string `json:"text"`
	Caption string `json:"caption"`
	Options []QuickReplyOption
}

func (qrt *QuickReplyText) MarshalJSON() ([]byte, error) {

	type TContent struct {
		Type    string `json:"type"`
		Header  string `json:"header"`
		Text    string `json:"text"`
		Caption string `json:"caption"`
	}

	tmpli := struct {
		MsgID   string `json:"msgid"`
		Content TContent
		Options []QuickReplyOption `json:"options"`
	}{
		MsgID: qrt.MsgID,
		Content: TContent{
			Type:    "text",
			Header:  qrt.Header,
			Text:    qrt.Text,
			Caption: qrt.Caption,
		},
		Options: qrt.Options,
	}

	return json.Marshal(tmpli)
}

type QuickReplyVideo struct {
	MsgID   string `json:"msgid"`
	URL     string `json:"url"`
	Text    string `json:"text"`
	Caption string `json:"caption"`
	Options []QuickReplyOption
}

func (qrv *QuickReplyVideo) MarshalJSON() ([]byte, error) {

	type TContent struct {
		Type    string `json:"type"`
		URL     string `json:"url"`
		Text    string `json:"text"`
		Caption string `json:"caption"`
	}

	tmpli := struct {
		MsgID   string `json:"msgid"`
		Content TContent
		Options []QuickReplyOption `json:"options"`
	}{
		MsgID: qrv.MsgID,
		Content: TContent{
			Type:    "video",
			URL:     qrv.URL,
			Text:    qrv.Text,
			Caption: qrv.Caption,
		},
		Options: qrv.Options,
	}

	return json.Marshal(tmpli)
}

type QuickReplyDocument struct {
	MsgID    string `json:"msgid"`
	URL      string `json:"url"`
	Text     string `json:"text"`
	Caption  string `json:"caption"`
	Filename string `json:"filename"`
	Options  []QuickReplyOption
}

func (qrd *QuickReplyDocument) MarshalJSON() ([]byte, error) {

	type TContent struct {
		Type     string `json:"type"`
		URL      string `json:"url"`
		Text     string `json:"text"`
		Caption  string `json:"caption"`
		Filename string `json:"filename"`
	}

	tmpli := struct {
		MsgID   string             `json:"msgid"`
		Content TContent           `json:"content"`
		Options []QuickReplyOption `json:"options"`
	}{
		MsgID: qrd.MsgID,
		Content: TContent{
			Type:     "document",
			URL:      qrd.URL,
			Text:     qrd.Text,
			Caption:  qrd.Caption,
			Filename: qrd.Filename,
		},
		Options: qrd.Options,
	}

	return json.Marshal(tmpli)
}
