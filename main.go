package main

import (
	"flag"
	"fmt"
	"time"

	v1 "github.com/jonasdacruz/1brc/v1"
	v2 "github.com/jonasdacruz/1brc/v2"
)

func main() {
	defer func(s time.Time) {
		fmt.Println("time taken: ", time.Since(s))
	}(time.Now())

	pv := flag.String("v", "v1", "print version")

	filePath := "./data/measurements.txt"

	flag.Parse()
	switch *pv {
	case "v1":
		fmt.Println("running v1")
		v1.ProcessFile(filePath)
	case "v2":
		fmt.Println("running v2")
		v2.ProcessFile(filePath)
	default:
		fmt.Println("version not found")
	}
}
