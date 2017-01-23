package slack_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"

	"github.com/chielkunkels/marvin"
	"github.com/chielkunkels/marvin/adapter/slack"
)

var testToken = "xoxb-1234567890-aAbBcCdDeEfFgGhHiIjJkKlL"

func TestOpen(t *testing.T) {
	tests := []struct {
		closeEarly  bool
		err         string
		parseURL    bool
		resRtmStart string
		rtmHit      bool
		rtmStartHit bool
		upgrade     bool
	}{
		{
			closeEarly: true,
			err:        slack.ErrHTTPStart.Error(),
		},
		{
			err:         "unexpected end of JSON input",
			resRtmStart: `{"foo":"bar"`,
			rtmStartHit: true,
		},
		{
			err:         "invalid_auth",
			resRtmStart: `{"ok":false,"error":"invalid_auth"}`,
			rtmStartHit: true,
		},
		{
			err:         "malformed ws or wss URL",
			resRtmStart: "{\"ok\":true}",
			rtmStartHit: true,
		},
		{
			err:         "websocket: bad handshake",
			parseURL:    true,
			resRtmStart: "{\"ok\":true,\"url\":\"%s\"}",
			rtmStartHit: true,
			rtmHit:      true,
		},
		{
			err:         "",
			parseURL:    true,
			resRtmStart: "{\"ok\":true,\"url\":\"%s\"}",
			rtmStartHit: true,
			rtmHit:      true,
			upgrade:     true,
		},
	}

	for i, test := range tests {
		var URL *url.URL
		rtmStartHit := false
		rtmHit := false

		h := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/rtm.start" {
				rtmStartHit = true

				res := test.resRtmStart
				if test.parseURL {
					URL.Scheme = "ws"
					res = fmt.Sprintf(res, URL.String()+"/rtm")
				}

				w.Write([]byte(res))
			}

			if r.URL.Path == "/rtm" {
				rtmHit = true

				if test.upgrade {
					upgrader := websocket.Upgrader{
						ReadBufferSize:  1024,
						WriteBufferSize: 1024,
					}

					conn, err := upgrader.Upgrade(w, r, nil)
					if err != nil {
						fmt.Println(err)
						return
					}
					defer conn.Close()
				}
			}
		}

		ts := httptest.NewServer(http.HandlerFunc(h))
		URL, _ = url.Parse(ts.URL)

		adapter := slack.NewAdapter(testToken)
		adapter.RtmStartEndpoint = URL.String() + "/rtm.start?token=%s"

		if test.closeEarly {
			ts.Close()
		}

		messages := make(chan *marvin.Message)
		err := adapter.Open(messages)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}

		if errStr != test.err {
			t.Errorf("%d: error: expected %+v, got %s\n", i, test.err, err)
		}

		if rtmStartHit != test.rtmStartHit {
			t.Errorf("%d: test was supposed to hit /rtm.start but did not (or vice versa)", i)
		}

		if rtmHit != test.rtmHit {
			t.Errorf("%d: test was supposed to hit /rtm but did not (or vice versa)", i)
		}

		if !test.closeEarly {
			ts.Close()
		}
	}
}