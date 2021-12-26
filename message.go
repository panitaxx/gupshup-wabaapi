package wabaapi

import (
	"encoding/json"
	"time"

	"github.com/ansel1/merry"
)

type InboundMessagePayload struct {
	ID      string      `json:"id"`
	Source  string      `json:"source"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (msg *InboundMessagePayload) UnmarshalJSON(data []byte) error {
	var tmp struct {
		ID      string          `json:"id"`
		Source  string          `json:"source"`
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
		Sender  Sender          `json:"sender"`
		Context *Context        `json:"context"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	msg.ID = tmp.ID
	msg.Source = tmp.Source
	msg.Type = tmp.Type

	switch msg.Type {
	case "text":
		tmpt := struct {
			Text string `json:"text"`
			Type string `json:"type"`
		}{}
		if err := json.Unmarshal(tmp.Payload, &tmpt); err != nil {
			return merry.Errorf("failed to parse text payload: %s", err)
		}
		if tmpt.Type == "button" {
			msg.Payload = InboundButtonText(tmpt.Text)
		}
		msg.Payload = InboundText(tmpt.Text)
	case "audio", "video", "image", "sticker", "file":
		var media InboundMedia
		if err := json.Unmarshal(tmp.Payload, &media); err != nil {
			return merry.Errorf("failed to parse media payload: %s", err)
		}
		media.Type = msg.Type
		msg.Payload = media
	case "location":
		var loc InboundLocation
		if err := json.Unmarshal(tmp.Payload, &loc); err != nil {
			return merry.Errorf("failed to parse location payload: %s", err)
		}
		msg.Payload = loc
	case "contact":
		var contact struct {
			Contacts []Contact `json:"contacts"`
		}
		if err := json.Unmarshal(tmp.Payload, &contact); err != nil {
			return merry.Errorf("failed to parse contact payload: %s", err)
		}
		msg.Payload = contact.Contacts
	case "list_reply":
		var list InboundListReply
		if err := json.Unmarshal(tmp.Payload, &list); err != nil {
			return merry.Errorf("failed to parse list_reply payload: %s", err)
		}
		msg.Payload = list
	case "button_reply":
		var btn InboundButtonReply
		if err := json.Unmarshal(tmp.Payload, &btn); err != nil {
			return merry.Errorf("failed to parse button_reply payload: %s", err)
		}
		msg.Payload = btn

	}

	return nil
}

type Sender struct {
	Phone       string `json:"phone"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	DialCode    string `json:"dial_code"`
}

type Context struct {
	ID   string `json:"id"`
	GsID string `json:"gsId"`
}

type InboundText string

type InboundButtonText string

type InboundMedia struct {
	Caption     string    `json:"caption"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	ContentType string    `json:"contentType"`
	URLExpiry   time.Time `json:"urlExpiry"`
	Type        string    `json:"-"`
}

func (media *InboundMedia) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Caption     string `json:"caption"`
		Name        string `json:"name"`
		URL         string `json:"url"`
		ContentType string `json:"contentType"`
		URLExpiry   int64  `json:"urlExpiry"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	media.Caption = tmp.Caption
	media.Name = tmp.Name
	media.URL = tmp.URL
	media.ContentType = tmp.ContentType
	media.URLExpiry = time.Now().Add(time.Hour)

	media.URLExpiry = time.Unix(tmp.URLExpiry/1000, 0)

	return nil
}

type InboundListReply struct {
	Title        string `json:"title"`
	ID           string `json:"id"`
	Reply        string `json:"reply"`
	PostbackText string `json:"postbackText"`
	Description  string `json:"description"`
}

type InboundButtonReply struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Reply string `json:"reply"`
}

type InboundLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Contact struct {
	Addresses []struct {
		City        string `json:"city"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		State       string `json:"state"`
		Street      string `json:"street"`
		Type        string `json:"type"`
		Zip         string `json:"zip"`
	} `json:"addresses"`
	Emails []struct {
		Email string `json:"email"`
		Type  string `json:"type"`
	} `json:"emails"`
	Ims  []interface{} `json:"ims"`
	Name struct {
		FirstName     string `json:"first_name"`
		FormattedName string `json:"formatted_name"`
		LastName      string `json:"last_name"`
	} `json:"name"`
	Org struct {
		Company string `json:"company"`
	} `json:"org"`
	Phones []struct {
		Phone string `json:"phone"`
		Type  string `json:"type"`
	} `json:"phones"`
	URLs []struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"urls"`
}
