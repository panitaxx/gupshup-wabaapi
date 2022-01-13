package wabaapi

import (
	"encoding/json"
	"time"

	"github.com/ansel1/merry"
)

type InboundMessage struct {
	App       string      `json:"app"`
	Timestamp time.Time   `json:"timestamp"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
}

func (m *InboundMessage) UnmarshalJSON(data []byte) error {
	type TMPMsg struct {
		App       string          `json:"app"`
		Timestamp int64           `json:"timestamp"`
		Type      string          `json:"type"`
		Payload   json.RawMessage `json:"payload"`
	}

	var tmp TMPMsg

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	t := time.Unix(tmp.Timestamp/1000, 0)

	*m = InboundMessage{
		App:       tmp.App,
		Timestamp: t,
		Type:      tmp.Type,
	}

	var payload interface{}
	switch tmp.Type {
	case "user-event":
		var tmpP UserEventPayload
		if err := json.Unmarshal(tmp.Payload, &tmpP); err != nil {
			return merry.Errorf("failed to parse user-event payload: %s", err)
		}
		payload = tmpP
	case "system-event":
		var tmpP SystemEventPayload
		if err := json.Unmarshal(tmp.Payload, &tmpP); err != nil {
			return merry.Errorf("failed to parse system-event payload: %s", err)
		}
		payload = tmpP
	case "account-event":

		var tmpP AccountEventPayload
		if err := json.Unmarshal(tmp.Payload, &tmpP); err != nil {
			return merry.Errorf("failed to parse account-event payload: %s", err)
		}
		payload = tmpP
	case "message-event":
		var tmpP MessageEventPayload
		if err := json.Unmarshal(tmp.Payload, &tmpP); err != nil {
			return merry.Errorf("failed to parse message-event payload: %s", err)
		}
		payload = tmpP
	case "message":
		var tmpP InboundMessagePayload
		if err := json.Unmarshal(tmp.Payload, &tmpP); err != nil {
			return merry.Errorf("failed to parse message payload: %s", err)
		}
		payload = tmpP

	}

	m.Payload = payload
	return nil
}

type UserEventPayload struct {
	Phone string `json:"phone"`
	Type  string `json:"type"`
}

type SystemEventPayload struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	ElementName    string `json:"elementName"`
	LanguageCode   string `json:"languageCode"`
	RejectedReason string `json:"rejectedReason"`
}

type AccountEventPayload struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

type MessageEventPayload struct {
	ID          string          `json:"id"`
	GSID        string          `json:"gsId"`
	Type        string          `json:"type"`
	Destination string          `json:"destination"`
	Payload     json.RawMessage `json:"payload"`
}

func (msgEvent *MessageEventPayload) GetError() error {
	if msgEvent.Type != "failed" {
		return nil
	}

	var payload struct {
		Code   int    `json:"code"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(msgEvent.Payload, &payload); err != nil {
		return merry.Errorf("failed to parse message-event payload: %s", err)
	}
	return merry.Errorf("message-event failed:[%d] %s", payload.Code, payload.Reason)
}
