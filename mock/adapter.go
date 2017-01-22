package mock

import "github.com/chielkunkels/marvin"

// Adapter represents a mock adapter
type Adapter struct {
	err      error
	messages chan<- *marvin.Message

	CloseCalled bool
	OpenCalled  bool
	ReplyCalled bool
	SendCalled  bool
}

// NewAdapter returns a new mock adapter
func NewAdapter() *Adapter {
	return &Adapter{}
}

// Close mocks an adapter closing the connection
func (a *Adapter) Close() error {
	a.CloseCalled = true
	return a.err
}

// Open mocks an adapter opening the connection
func (a *Adapter) Open(messages chan<- *marvin.Message) error {
	a.messages = messages
	a.OpenCalled = true
	return a.err
}

// PushMessage pushes a new message into the messages channel
func (a *Adapter) PushMessage(m *marvin.Message) {
	a.messages <- m
}

// Reply sends a reply directed at the user sending the request
func (a *Adapter) Reply(m *marvin.Message, text string) error {
	a.ReplyCalled = true
	return a.err
}

// Send sends a message in the channel the request originated from
func (a *Adapter) Send(m *marvin.Message, text string) error {
	a.SendCalled = true
	return a.err
}

// SetError sets an error
func (a *Adapter) SetError(err error) {
	a.err = err
}
