package cv

import (
	"image/png"
	"os"
	"testing"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

func TestGetEdges(t *testing.T) {
	img := util.GetRandomImage()

	edges, err := GetEdges(img)
	if err != nil {
		t.Error(err)
	}

	inputFile, err := os.Create("input.png")
	if err != nil {
		t.Error(err)
	}
	defer inputFile.Close()

	err = png.Encode(inputFile, img)
	if err != nil {
		t.Error(err)
	}

	edgesFile, err := os.Create("edges.png")
	if err != nil {
		t.Error(err)
	}
	defer edgesFile.Close()

	err = png.Encode(edgesFile, edges)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPalette(t *testing.T) {
	img := util.GetRandomImage()

	palette, err := GetPalette(img, 64, false)
	if err != nil {
		t.Error(err)
	}

	paletteImg := GetPaletteImage(palette)

	inputFile, err := os.Create("input.png")
	if err != nil {
		t.Error(err)
	}
	defer inputFile.Close()

	err = png.Encode(inputFile, img)
	if err != nil {
		t.Error(err)
	}

	paletteFile, err := os.Create("palette.png")
	if err != nil {
		t.Error(err)
	}
	defer paletteFile.Close()

	err = png.Encode(paletteFile, paletteImg)
	if err != nil {
		t.Error(err)
	}
}
