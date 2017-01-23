package slack

// Slack errors
const (
	ErrHTTPStart = Error("failed to make call to rtm.start")
)

// Error describes a Slack error
type Error string

// Error returns the error
func (e Error) Error() string {
	return string(e)
}
