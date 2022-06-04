package main

import (
	"fmt"
	"log"
)

func main() {
	driver := NewAWSS3Driver()
	if err := driver.Start(); err != nil {
		log.Fatalln(err)
	}

	buckets, err := driver.ListBuckets()

	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to list buckets: %s", err.Error()))
	}

	for _, buck := range buckets.Buckets {
		println(buck.GoString(), buck.CreationDate.GoString())
	}
}
