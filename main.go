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

	termbox.PollEvent()
}
