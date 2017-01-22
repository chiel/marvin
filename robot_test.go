package marvin_test

import (
	"errors"
	"testing"

	"github.com/chielkunkels/marvin"
	"github.com/chielkunkels/marvin/mock"
)

func TestNewRobot(t *testing.T) {
	adapter := mock.NewAdapter()

	_, err := marvin.NewRobot("mar[vin", adapter)
	if err == nil {
		t.Error("NewRobot should have failed with name `mar[vin`")
	}

	_, err = marvin.NewRobot("marvin", adapter)
	if err != nil {
		t.Error("NewRobot should not have failed with name `marvin`")
	}
}

func Test_receiveMessages(t *testing.T) {
	tests := []struct {
		Called  bool
		Direct  bool
		Pattern string
		Text    string
	}{
		{true, true, "^test", "@marvin test"},
		{false, true, "^test", "test"},
		{true, false, "test", "@marvin test"},
		{false, false, "test", "thing"},
	}

	for _, test := range tests {
		adapter := mock.NewAdapter()
		robot, _ := marvin.NewRobot("marvin", adapter)
		robot.Open()

		m := &marvin.Message{
			Channel: &marvin.Channel{ID: "1234", Name: "general"},
			User:    &marvin.User{ID: "4321", Name: "someperson"},
			Text:    test.Text,
		}

		called := false
		callback := func(r *marvin.Request) {
			called = true
		}

		if test.Direct {
			robot.Respond(test.Pattern, callback)
		} else {
			robot.Hear(test.Pattern, callback)
		}

		adapter.PushMessage(m)

		if called != test.Called {
			t.Error("Should not have been called")
		}
	}
}

func TestClose(t *testing.T) {
	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)
	if err := robot.Close(); err != nil {
		t.Error("Close should not have returned an error")
	}

	if !adapter.CloseCalled {
		t.Error("Close was not called on adapter")
	}

	adapter = mock.NewAdapter()
	adapter.SetError(errors.New("oh noes"))
	robot, _ = marvin.NewRobot("marvin", adapter)
	if err := robot.Close(); err == nil {
		t.Error("Close should have returned an error")
	}
}

func TestHear(t *testing.T) {
	cb := func(*marvin.Request) {}

	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)
	if err := robot.Hear("test", cb); err != nil {
		t.Error("Hear should not have returned an error")
	}

	if err := robot.Hear("^te[st", cb); err == nil {
		t.Error("Hear should have returned an error")
	}
}

func TestOpen(t *testing.T) {
	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)
	if err := robot.Open(); err != nil {
		t.Error("Open should not have returned an error")
	}

	if !adapter.OpenCalled {
		t.Error("Open was not called on adapter")
	}

	adapter = mock.NewAdapter()
	adapter.SetError(errors.New("oh noes"))
	robot, _ = marvin.NewRobot("marvin", adapter)
	if err := robot.Open(); err == nil {
		t.Error("Open should have returned an error")
	}
}

func TestRegisterPlugin(t *testing.T) {
	pluginCalled := false

	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)
	robot.RegisterPlugin(func(r *marvin.Robot) {
		pluginCalled = true

		if r != robot {
			t.Error("Did not get passed the correct robot")
		}
	})

	if !pluginCalled {
		t.Error("Plugin did not get called")
	}
}

func TestRespond(t *testing.T) {
	cb := func(*marvin.Request) {}

	adapter := mock.NewAdapter()
	robot, _ := marvin.NewRobot("marvin", adapter)
	if err := robot.Respond("test", cb); err != nil {
		t.Error("Respond should not have returned an error")
	}

	if err := robot.Respond("^te[st", cb); err == nil {
		t.Error("Respond should have returned an error")
	}
}
