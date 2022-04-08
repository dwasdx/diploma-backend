package dispatchers

import (
	"fmt"
	"log"
	"shopingList/pkg/listeners"
)

type SimpleEventDispatcher struct {
	events    map[string][]listeners.Event
	listeners map[string][]listeners.Listener
}

func NewSimpleEventDispatcher() SimpleEventDispatcher {
	return SimpleEventDispatcher{listeners: map[string][]listeners.Listener{}}
}

func (s *SimpleEventDispatcher) RegisterListener(event listeners.Event, listener listeners.Listener) {
	eventType := s.getEventType(event)

	s.listeners[eventType] = append(s.listeners[eventType], listener)
}

func (s *SimpleEventDispatcher) Dispatch(event listeners.Event) {
	typeEvent := s.getEventType(event)

	if s.events == nil {
		s.events = make(map[string][]listeners.Event, 0)
	}

	if s.events[typeEvent] == nil {
		s.events[typeEvent] = make([]listeners.Event, 0)
	}

	s.events[typeEvent] = append(s.events[typeEvent], event)
}

func (s *SimpleEventDispatcher) getEventType(event interface{}) string {
	return fmt.Sprintf("%T", event)
}

func (s *SimpleEventDispatcher) ReleaseAll() {
	log.Println("SimpleEventDispatcher->ReleaseAll()")

	for eventType, events := range s.events {
		if s.listeners[eventType] != nil {
			for ind, listener := range s.listeners[eventType] {
				for _, event := range events {
					listener.Handle(event)

					s.events[eventType] = append(s.events[eventType][:ind], s.events[eventType][ind+1:]...)
				}
			}
		}
	}
}
