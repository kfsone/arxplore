package main


import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"io"
	"os"
)


// Base worker that consumes requests from the queue, opens the requested
// file and forwards it to the relevant listing function for that file type.
//
func worker(id int) {
	defer wg.Done()
	for request := range workRequest {
		file, err := os.Open(request.Path)
		if err != nil {
			return
		}
		func() {
			defer file.Close()
			request.Function(request.Path, file)
		}()
	}
}

// List the files from a tar archive
func untar(path string, stream io.Reader) {
	reader := tar.NewReader(stream)
	for {
		header, err := reader.Next()
		if err != nil {
			break
		}
		workResult <- WorkResult{path, header.Name}
	}
}

// decode a bzip2 encoded file and forward to the untar function
func unbzip2(path string, file io.Reader) {
	untar(path, bzip2.NewReader(file))
}

// decode a gzip encoded file and forward to the untar function
func ungzip(path string, file io.Reader) {
	stream, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer stream.Close()
	untar(path, stream)
}

// list the files in a .zip file
func unzip(path string, file io.Reader) {
	zipfile, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer zipfile.Close()
	for _, f := range zipfile.File {
		workResult <- WorkResult{path, f.Name}
	}
}

