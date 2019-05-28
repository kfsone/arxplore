package main


import "io"


type WorkerFn func(string, io.Reader)

type WorkRequest struct {
	Function WorkerFn
	Path     string
	Size     int64
}

type WorkResult struct {
	Archive string
	File    string
}

