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
	defer termbox.Close()

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

	for {
		ev := <-eventQueue
		switch {
		case ev.Key == termbox.KeyArrowUp:
			view.Up()
		case ev.Key == termbox.KeyArrowDown:
			view.Down()
		case ev.Key == termbox.KeyEnter:
			view.Dive()
		case ev.Key == termbox.KeyEsc:
			termbox.SetCursor(0, 0)
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			os.Exit(0)
		}
	}
}
