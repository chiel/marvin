package marvin

import "regexp"

// Adapter describes the interface an adapter should implement.
type Adapter interface {
	Close() error
	Open() error
	Reply(*Message, string) error
	Send(*Message, string) error
}

// Channel describes a channel.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IsDM bool
}

// Listener describes a listener.
type Listener struct {
	callback ListenerCallback
	direct   bool
	regex    *regexp.Regexp
}

// ListenerCallback describes the signature of a listener callback.
type ListenerCallback func(*Request)

// Message describes a message.
type Message struct {
	Channel *Channel
	User    *User
	Text    string
}

// User describes a user.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
