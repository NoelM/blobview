package main

import (
	"fmt"
	"log"
	"time"
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
		println(*buck.Name, buck.CreationDate.Format(time.RFC1123Z))
	}

	objects, err := driver.ListObjects(*buckets.Buckets[0].Name, "")
	if err != nil {
		log.Fatalln(fmt.Sprintf("Unable to list objects: %s", err.Error()))
	}

	for _, obj := range objects.CommonPrefixes {
		println(*obj.Prefix)
	}
}
