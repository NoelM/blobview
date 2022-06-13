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

	for {
		ev := <-eventQueue
		switch ev.Type {
		case termbox.EventKey:
			ch := ev.Ch
			// handle characters first, see termbox.Event
			if ch != 0 {
				switch ch {
				case 'd':
					view.Download()
				case 'j':
					view.Down()
				case 'k':
					view.Up()
				case 'h':
					view.Back()
				case 'l':
					view.Dive()
				case 'q':
					termbox.Close()
					os.Exit(0)
				}
			}

			key := ev.Key
			switch key {
			case termbox.KeyArrowUp:
				view.Up()
			case termbox.KeyArrowDown:
				view.Down()
			case termbox.KeyEnter:
				view.Dive()
			case termbox.KeyBackspace2:
				view.Back()
			case termbox.KeyEsc:
				termbox.Close()
				os.Exit(0)
			}
		}

	}
}
