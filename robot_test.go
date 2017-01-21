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
