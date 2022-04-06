package dispatchers

import (
	"shopingList/pkg/listeners"
)

type EventDispatcher interface {
	RegisterListener(event listeners.Event, listener listeners.Listener)
	Dispatch(event listeners.Event)
	ReleaseAll()
}
