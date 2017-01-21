package marvin

// Adapter describes the interface an adapter should implement.
type Adapter interface {
	Close() error
	Open() error
	Reply(*Message, string) error
	Send(*Message, string) error
}

// Channel describes a channel
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IsDM bool
}

// Message describes a message
type Message struct {
	Channel *Channel
	User    *User
	Text    string
}

// User describes a user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
