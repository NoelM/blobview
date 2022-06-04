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

	driver := NewAWSS3Driver()
	if err := driver.Start(); err != nil {
		log.Fatalln(err)
	}

	buckets, err := driver.ListBuckets()
	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to list buckets: %s", err.Error()))
	}
	view := NewView()
	view.PrintObjectList(buckets)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		ev := <-eventQueue
		switch ev.Key {
		case termbox.KeyArrowUp:
			view.Up()
		case termbox.KeyArrowDown:
			view.Down()
		case 'q':
			os.Exit(0)
		}
	}
}
