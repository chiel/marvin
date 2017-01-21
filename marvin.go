package marvin

// Adapter describes the interface an adapter should implement.
type Adapter interface {
	Close() error
	Open() error
}
