package marvin

// Request describes an incoming request.
type Request struct {
	Message *Message
	Query   []string
	robot   *Robot
}

// NewRequest creates a new request and return a pointer to it.
func NewRequest(robot *Robot, message *Message, query []string) *Request {
	return &Request{
		Message: message,
		Query:   query,
		robot:   robot,
	}
}

// Reply sends a reply to the user sending the request.
func (r *Request) Reply(text string) {
	r.robot.adapter.Reply(r.Message, text)
}

// Send sends a message to the channel the request originated from.
func (r *Request) Send(text string) {
	r.robot.adapter.Send(r.Message, text)
}
