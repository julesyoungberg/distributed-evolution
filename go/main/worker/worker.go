package main

import (
	"encoding/gob"
	"image"

	"github.com/rickyfitts/distributed-evolution/worker"
)

func main() {
	// TODO fix??
	gob.Register(image.YCbCr{})
	gob.Register(worker.Triangle{})
	worker.Run()
}
