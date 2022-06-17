package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to init termbox: %s", err.Error()))
	}

	view := NewObjectListView()
	if err := view.Start(); err != nil {
		log.Fatalln(fmt.Sprintf("Unable to start view: %s", err.Error()))
	}

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	handler := NewEventHandler(eventQueue, view)

	handler.RegisterAction(NewViewAction("download", func(view *ObjectListView) { view.Download() }, 'd'))
	handler.RegisterAction(NewViewAction("up", func(view *ObjectListView) { view.Up() }, 'k', termbox.KeyArrowUp))
	handler.RegisterAction(NewViewAction("down", func(view *ObjectListView) { view.Down() }, 'j', termbox.KeyArrowDown))
	handler.RegisterAction(NewViewAction("back", func(view *ObjectListView) { view.Back() }, 'h', termbox.KeyBackspace2))
	handler.RegisterAction(NewViewAction("dive", func(view *ObjectListView) { view.Dive() }, 'l', termbox.KeyEnter))
	handler.RegisterAction(NewViewAction("close", func(view *ObjectListView) { termbox.Close(); os.Exit(0) }, 'q', termbox.KeyEsc))

	handler.Start()
}
