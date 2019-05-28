package main

import (
	"log"
	"os"
	"strings"
)

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("Skipping '%s': %s\n", path, err)
		return nil
	}

	if info.IsDir() {
		return nil
	}

	var workerFn WorkerFn
	switch {
	case strings.HasSuffix(path, ".zip"):
		{
			workerFn = unzip
		}
	case strings.HasSuffix(path, ".tbz2"), strings.HasSuffix(path, ".tar.bz2"):
		{
			workerFn = unbzip2
		}
	case strings.HasSuffix(path, ".tgz"), strings.HasSuffix(path, ".tar.gz"):
		{
			workerFn = ungzip
		}
	case strings.HasSuffix(path, ".tar"):
		{
			workerFn = untar
		}
	}

	if workerFn != nil {
		workRequest <- WorkRequest{workerFn, path, info.Size()}
	}

	return nil
}

