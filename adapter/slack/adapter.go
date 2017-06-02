package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/chielkunkels/marvin"
)

// Adapter describes a slack adapter.
type Adapter struct {
	channelsByID     map[string]*marvin.Channel
	channelsByName   map[string]*marvin.Channel
	counter          int64
	RtmStartEndpoint string
	self             marvin.User
	token            string
	usersByID        map[string]*marvin.User
	usersByName      map[string]*marvin.User
	ws               *websocket.Conn
}

// NewAdapter creates a new slack adapter.
func NewAdapter(token string) *Adapter {
	return &Adapter{
		channelsByID:     map[string]*marvin.Channel{},
		channelsByName:   map[string]*marvin.Channel{},
		RtmStartEndpoint: "https://slack.com/api/rtm.start?token=%s",
		token:            token,
		usersByID:        map[string]*marvin.User{},
		usersByName:      map[string]*marvin.User{},
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

// cacheChannels takes all the slack channels from
// the rtm.start response and caches them in memory.
func (a *Adapter) cacheChannels(channels []marvin.Channel) {
	for _, c := range channels {
		c := c
		c.IsDM = strings.HasPrefix(c.ID, "D")
		a.channelsByID[c.ID] = &c
		a.channelsByName[c.Name] = &c
	}
}

// cacheUsers takes all the users from the rtm.start
// response and caches them in memory.
func (a *Adapter) cacheUsers(users []marvin.User) {
	for _, u := range users {
		u := u
		a.usersByID[u.ID] = &u
		a.usersByName[u.Name] = &u
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

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var res rtmStart
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	if !res.Ok {
		return errors.New(res.Err)
	}

	a.self = res.Self
	a.cacheChannels(res.Channels)
	a.cacheChannels(res.Groups)
	a.cacheChannels(res.IMs)
	a.cacheUsers(res.Users)

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
