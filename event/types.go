package event

type Fetcher interface {
	// fetct events
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	// process event
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
