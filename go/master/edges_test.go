package master

import (
	"image/png"
	"os"
	"testing"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

func TestGetEdges(t *testing.T) {
	img := util.GetRandomImage()

	edges, err := getEdges(img)
	if err != nil {
		t.Error(err)
	}

	inputFile, err := os.Create("input.png")
	if err != nil {
		t.Error(err)
	}

	err = png.Encode(inputFile, img)
	if err != nil {
		t.Error(err)
	}

	edgesFile, err := os.Create("edges.png")
	if err != nil {
		t.Error(err)
	}

	err = png.Encode(edgesFile, edges)
	if err != nil {
		t.Error(err)
	}
}
