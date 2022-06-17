package main

import (
	"github.com/nsf/termbox-go"
	"os"
)

type ViewAction struct {
	name          string
	triggerEvents []termbox.Event
	cb            func(view *ObjectListView)
}

type EventHandler struct {
	eventQueue   chan termbox.Event
	event2Action map[termbox.Event]ViewAction
	view         *ObjectListView
}

func (EH *EventHandler) RegisterAction(action ViewAction) {
	for _, ev := range action.triggerEvents {
		if _, ok := EH.event2Action[ev]; ok {
			panic("event already linked to another action")
		}
		EH.event2Action[ev] = action
	}
}

func NewEventHandler(ev chan termbox.Event, view *ObjectListView) *EventHandler {
	m := make(map[termbox.Event]ViewAction)
	return &EventHandler{eventQueue: ev, event2Action: m, view: view}
}

func (EH *EventHandler) Start() {
	for {
		ev := <-EH.eventQueue
		// reset n-bytes of event to have map access works
		ev.N = 0

		if action, ok := EH.event2Action[ev]; ok {
			action.cb(EH.view)
		}
	}
}

func NewChEvent(ch rune) termbox.Event {
	return termbox.Event{
		Type: termbox.EventKey,
		Ch:   ch,
	}
}

func NewKeyEvent(key termbox.Key) termbox.Event {
	return termbox.Event{
		Type: termbox.EventKey,
		Key:  key,
	}
}

func NewViewAction(name string, cb func(view *ObjectListView), evs ...interface{}) ViewAction {
	v := ViewAction{name: name, cb: cb}
	for _, ev := range evs {
		switch evValue := ev.(type) {
		case rune:
			v.triggerEvents = append(v.triggerEvents, NewChEvent(evValue))
		case termbox.Key:
			v.triggerEvents = append(v.triggerEvents, NewKeyEvent(evValue))
		default:
			panic("unable to build action ViewAction")
		}
	}
	return v
}

func CreateActions() []ViewAction {
	actions := make([]ViewAction, 0)
	actions = append(actions, NewViewAction("download", func(view *ObjectListView) { view.Download() }, 'd'))
	actions = append(actions, NewViewAction("up", func(view *ObjectListView) { view.Up() }, 'k', termbox.KeyArrowUp))
	actions = append(actions, NewViewAction("down", func(view *ObjectListView) { view.Down() }, 'j', termbox.KeyArrowDown))
	actions = append(actions, NewViewAction("back", func(view *ObjectListView) { view.Back() }, 'h', termbox.KeyBackspace2))
	actions = append(actions, NewViewAction("dive", func(view *ObjectListView) { view.Dive() }, 'l', termbox.KeyEnter))
	actions = append(actions, NewViewAction("close", func(view *ObjectListView) { termbox.Close(); os.Exit(0) }, 'q', termbox.KeyEsc))

	return actions
}

func NewViewActionMap() map[termbox.Event]ViewAction {
	actions := CreateActions()
	m := make(map[termbox.Event]ViewAction)
	for _, action := range actions {
		for _, ev := range action.triggerEvents {
			m[ev] = action
		}
	}
	return m
}
