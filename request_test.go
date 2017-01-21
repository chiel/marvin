package marvin_test

import (
	"testing"

	"github.com/chielkunkels/marvin"
	"github.com/chielkunkels/marvin/mock"
)

func TestNewRequest(t *testing.T) {
	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)

	m := &marvin.Message{
		Channel: &marvin.Channel{ID: "1234", Name: "general"},
		User:    &marvin.User{ID: "4321", Name: "someperson"},
		Text:    "Testing!",
	}

	query := []string{"thing"}

	request := marvin.NewRequest(robot, m, query)
	if request.Message != m || request.Query[0] != query[0] {
		t.Error("Request was not instantiated properly")
	}
}

func TestReply(t *testing.T) {
	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)

	m := &marvin.Message{
		Channel: &marvin.Channel{ID: "1234", Name: "general"},
		User:    &marvin.User{ID: "4321", Name: "someperson"},
		Text:    "Testing!",
	}

	request := marvin.NewRequest(robot, m, []string{})
	request.Reply("stuff and things")

	if !adapter.ReplyCalled {
		t.Error("Reply was not called on the adapter")
	}
}

func TestSend(t *testing.T) {
	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)

	m := &marvin.Message{
		Channel: &marvin.Channel{ID: "1234", Name: "general"},
		User:    &marvin.User{ID: "4321", Name: "someperson"},
		Text:    "Testing!",
	}

	request := marvin.NewRequest(robot, m, []string{})
	request.Send("stuff and things")

	if !adapter.SendCalled {
		t.Error("Reply was not called on the adapter")
	}
}
