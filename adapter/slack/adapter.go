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

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return err
	}

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
