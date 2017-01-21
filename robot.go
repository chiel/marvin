package marvin

import "regexp"

// Robot describes a robot.
type Robot struct {
	adapter   Adapter
	listeners []*Listener
	name      string
	nameRegex *regexp.Regexp
}

// NewRobot creates a new robot and returns a pointer to it.
func NewRobot(name string, adapter Adapter) (*Robot, error) {
	nameRegex, err := regexp.Compile(`^@?` + name + `\:?\s+`)
	if err != nil {
		return nil, err
	}

	robot := &Robot{
		adapter:   adapter,
		name:      name,
		nameRegex: nameRegex,
	}

	return robot, nil
}

// Close disconnects the robot's adapter.
func (r *Robot) Close() error {
	return r.adapter.Close()
}

// Hear creates a listener for messages that are not necessarily directed at the robot.
func (r *Robot) Hear(pattern string, callback ListenerCallback) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	listener := &Listener{callback: callback, direct: false, regex: regex}
	r.listeners = append(r.listeners, listener)
	return nil
}

// Open connects the robot through the adapter.
func (r *Robot) Open() error {
	return r.adapter.Open()
}

// RegisterPlugin registers the given plugin.
func (r *Robot) RegisterPlugin(plugin func(*Robot)) {
	plugin(r)
}
