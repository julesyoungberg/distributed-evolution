package master

import (
	"fmt"
	"testing"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

func TestGetPalette(t *testing.T) {
	img := util.GetRandomImage()
	palette, err := getPalette(img, 8)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	for _, c := range palette {
		fmt.Printf("%v", c)
	}
}
