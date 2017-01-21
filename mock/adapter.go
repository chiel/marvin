package mock

// Adapter represents a mock adapter
type Adapter struct {
	err error

	CloseCalled bool
	OpenCalled  bool
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
func (a *Adapter) Open() error {
	a.OpenCalled = true
	return a.err
}

// SetError sets an error
func (a *Adapter) SetError(err error) {
	a.err = err
}
