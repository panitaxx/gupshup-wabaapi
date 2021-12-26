package wabaapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleOutboundMessage() {
	om := &OutboundMessage{
		Channel:     "whatsapp",
		Destination: "+1234567890",
		Source:      "+15555555555",
		SourceName:  "Our Company",
	}

	values, err := om.Text("Hello World")
	if err != nil {
		panic(err)
	}

	//Send the message
	_, err = http.PostForm("https://gupshupurl", values)
}

func TestListMessageMarshal(t *testing.T) {
	msg := ListMessage{
		Title:        "test_title",
		Body:         "test_body",
		MsgID:        "test_msgid",
		GlobalButton: "test_global_button",
		Items: []ListItem{
			{Title: "test1", Options: []ListItemOption{
				{Title: "test1_1", Description: "test1_1_desc", PostbackText: "test1_1_postback"},
				{Title: "test1_2", Description: "test1_2_desc", PostbackText: "test1_2_postback"},
			}},
		},
	}
	val, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"title":"test_title","body":"test_body","msgid":"test_msgid","type":"list","globalButtons":[{"type":"text","title":"test_global_button"}],"items":[
		{"title":"test1","options":[{"type":"text","title":"test1_1","description":"test1_1_desc","postbackText":"test1_1_postback"},{"type":"text","title":"test1_2","description":"test1_2_desc","postbackText":"test1_2_postback"}]}
	] }`, string(val))
}
