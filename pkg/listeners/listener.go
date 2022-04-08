package listeners

type Listener interface {
	Handle(event interface{}) error
}

type Event interface {
	GetEventType() string
}
