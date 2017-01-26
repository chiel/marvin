package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/chielkunkels/marvin"
)

// Adapter describes a slack adapter.
type Adapter struct {
	counter          int64
	RtmStartEndpoint string
	token            string
	ws               *websocket.Conn
}

// NewAdapter creates a new slack adapter.
func NewAdapter(token string) *Adapter {
	return &Adapter{
		RtmStartEndpoint: "https://slack.com/api/rtm.start?token=%s",
		token:            token,
	}
}

// sendMessage sends a message to slack's rtm api.
func (a *Adapter) sendMessage(m *marvin.Message, text string) error {
	a.counter++

	rm := &message{
		ID:      a.counter,
		Channel: m.Channel.ID,
		Text:    text,
		Type:    "message",
	}

	return a.ws.WriteJSON(rm)
}

// Close disconnects the adapter from slack's RTM api.
func (a *Adapter) Close() error {
	if a.ws != nil {
		a.ws.Close()
	}

	a.ws = nil
	return nil
}

// Open authenticates and connects to slack's RTM api.
func (a *Adapter) Open(messages chan<- *marvin.Message) error {
	url := fmt.Sprintf(a.RtmStartEndpoint, a.token)

	resp, err := http.Get(url)
	if err != nil {
		return ErrHTTPStart
	}

	body, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 1024))
	resp.Body.Close()

	var res rtmStart
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	if !res.Ok {
		return errors.New(res.Err)
	}

	ws, _, err := websocket.DefaultDialer.Dial(res.URL, nil)
	if err != nil {
		return err
	}

	a.ws = ws

	return nil
}

// Reply sends a reply to the user sending the request.
func (a *Adapter) Reply(m *marvin.Message, text string) error {
	if !m.Channel.IsDM {
		text = "@" + m.User.Name + " " + text
	}

	return a.sendMessage(m, text)
}

// Send sends some text back to the channel the message originated from.
func (a *Adapter) Send(m *marvin.Message, text string) error {
	return a.sendMessage(m, text)
}
