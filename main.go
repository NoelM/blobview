package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
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

	actionsMap := NewViewActionMap()
	for {
		ev := <-eventQueue
		// reset n-bytes of event to have map access works
		ev.N = 0

		if action, ok := actionsMap[ev]; ok {
			action.cb(view)
		}
	}
}
