package marvin

import (
	"net/http"
	"regexp"

	"github.com/pressly/chi"
)

// Robot describes a robot.
type Robot struct {
	adapter   Adapter
	address   string
	listeners []*Listener
	name      string
	nameRegex *regexp.Regexp
	plugins   []func(*Robot)
	Router    *chi.Mux
}

// NewRobot creates a new robot and returns a pointer to it.
func NewRobot(name string, adapter Adapter, address string) (*Robot, error) {
	nameRegex, err := regexp.Compile(`^@?` + name + `\:?\s+`)
	if err != nil {
		return nil, err
	}

	robot := &Robot{
		adapter:   adapter,
		address:   address,
		name:      name,
		nameRegex: nameRegex,
		plugins:   []func(*Robot){},
		Router:    chi.NewRouter(),
	}

	return robot, nil
}

// createListener adds a new listener.
func (r *Robot) createListener(pattern string, callback ListenerCallback, direct bool) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	listener := &Listener{callback: callback, direct: direct, regex: regex}
	r.listeners = append(r.listeners, listener)
	return nil
}

// receiveMessages listens for messages on the given channel.
func (r *Robot) receiveMessages(messages <-chan *Message) {
	for m := range messages {
		for _, listener := range r.listeners {
			if listener.direct && !r.nameRegex.MatchString(m.Text) {
				continue
			}

			text := m.Text
			if listener.direct {
				text = r.nameRegex.ReplaceAllString(m.Text, "")
			}

			matches := listener.regex.FindStringSubmatch(text)
			if matches == nil {
				continue
			}

			listener.callback(NewRequest(r, m, matches[1:]))
		}
	}
}

// Close disconnects the robot's adapter.
func (r *Robot) Close() error {
	return r.adapter.Close()
}

// Hear creates a listener for messages that are not necessarily directed at the robot.
func (r *Robot) Hear(pattern string, callback ListenerCallback) error {
	return r.createListener(pattern, callback, false)
}

// Open connects the robot through the adapter.
func (r *Robot) Open() error {
	messages := make(chan *Message)
	go r.receiveMessages(messages)

	go func() { http.ListenAndServe(r.address, r.Router) }()

	if err := r.adapter.Open(messages); err != nil {
		return err
	}

	for _, plugin := range r.plugins {
		plugin(r)
	}

	return nil
}

// RegisterPlugin registers the given plugin.
func (r *Robot) RegisterPlugin(plugin func(*Robot)) {
	r.plugins = append(r.plugins, plugin)
}

// Respond creates a listener for messages directed at the robot.
func (r *Robot) Respond(pattern string, callback ListenerCallback) error {
	return r.createListener(pattern, callback, true)
}

// Send sends text to a channel.
func (r *Robot) Send(channel string, text string) error {
	return r.adapter.SendMessage(channel, text)
}
