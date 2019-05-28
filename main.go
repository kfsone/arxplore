package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var workRequest chan WorkRequest
var workResult chan WorkResult
var wg sync.WaitGroup

func main() {
	var topFolder string
	var numWorkers int

	flag.StringVar(&topFolder, "top", ".", "Folder to start from")
	flag.IntVar(&numWorkers, "j", 1, "Number of concurrent workers to use")
	flag.Parse()
	if len(flag.Args()) > 0 {
		fmt.Println("ERROR: Unrecognized argument:", flag.Args()[0])
		fmt.Printf("Usage: %s [-top PATH] [-j N]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if topFolder == "" {
		fmt.Println("ERROR: No 'top' folder specified.")
		os.Exit(1)
	}
	if numWorkers < 1 {
		fmt.Println("ERROR: Invalid number of workers.")
		os.Exit(1)
	}

	log.Println("# Exploring:", topFolder, "with", numWorkers, "workers")

	workRequest = make(chan WorkRequest, 1024)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(int(i))
	}

	workResult = make(chan WorkResult, 4096)

	go func() {
		defer close(workRequest)
		err := filepath.Walk(topFolder, walkFn)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for result := range workResult {
			fmt.Printf("\"%s\": \"%s\"\n", result.Archive, result.File)
		}
	}()

	wg.Wait()

	close(workResult)
}
