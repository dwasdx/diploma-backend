package sync

import "shopingList/pkg/events"

type EventCollection struct {
	shareEvents []events.ShareListEvent
	goodEvents  []events.GoodsChangeEvent
}

func (s *EventCollection) AddShareEvent(event events.ShareListEvent) {
	s.shareEvents = append(s.shareEvents, event)
}

func (s *EventCollection) AddGoodEvent(event events.GoodsChangeEvent) {
	s.goodEvents = append(s.goodEvents, event)
}

func (s *EventCollection) GetShareEvents() []events.ShareListEvent {
	return s.shareEvents
}

func (s *EventCollection) GetGoodEvents() []events.GoodsChangeEvent {
	return s.goodEvents
}
