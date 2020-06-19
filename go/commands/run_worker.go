package main

import (
	"encoding/gob"
	"image"
	"image/color"

	"github.com/rickyfitts/distributed-evolution/worker"
)

func main() {
	// TODO fix??
	gob.Register(image.YCbCr{})
	gob.Register(worker.Triangle{})
	gob.Register(worker.Shapes{})
	gob.Register(color.RGBA{})
	worker.Run()
}
