package main

import (
	"image"

	"github.com/rickyfitts/distributed-evolution/util"
)

type Master struct {
	targetImage       image.Image
	targetImageBase64 string
}

func main() {
	m := new(Master)

	m.targetImage = util.GetRandomImage()
	m.targetImageBase64 = util.Base64EncodeImage(m.targetImage)

	go m.httpServer()

	m.rpcServer()
}
