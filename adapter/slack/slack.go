package slack

import "github.com/chielkunkels/marvin"

// rtmStart describes the structure of the rtm.start response
type rtmStart struct {
	Channels []marvin.Channel `json:"channels"`
	Err      string           `json:"error"`
	Groups   []marvin.Channel `json:"groups"`
	IMs      []marvin.Channel `json:"ims"`
	Ok       bool             `json:"ok"`
	Self     marvin.User      `json:"self"`
	URL      string           `json:"url"`
	Users    []marvin.User    `json:"users"`
}
