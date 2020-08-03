package master

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

// gets the color palette of an image with opencv k means clustering
func getPalette(img image.Image, k int) ([]color.RGBA, error) {
	log.Printf("[task-generator] getting palette from image")

	// convert input image to opencv matrix
	src, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return []color.RGBA{}, fmt.Errorf("getting palette: %v", err)
	}
	defer src.Close()

	data := src.Reshape(1, src.Total())
	data.ConvertTo(&data, gocv.MatTypeCV32F)

	labels := gocv.NewMat()
	criteria := gocv.NewTermCriteria(gocv.MaxIter+gocv.EPS, 10, 1.0)
	centers := gocv.NewMat()

	gocv.KMeans(data, k, &labels, criteria, 1, gocv.KMeansPPCenters, &centers)

	centers.ConvertTo(&centers, gocv.MatTypeCV8U)

	// collect palette colors
	palette := make([]color.RGBA, centers.Rows())
	for i := range palette {
		palette[i] = color.RGBA{
			R: uint8(centers.GetIntAt(i, 0)),
			G: uint8(centers.GetIntAt(i, 1)),
			B: uint8(centers.GetIntAt(i, 2)),
			A: 255,
		}
	}

	return palette, err
}

func (m *Master) preparePalette() {
	m.mu.Lock()
	img := m.TargetImage.Image
	nColors := m.Job.NumColors
	m.mu.Unlock()

	palette, err := getPalette(img, nColors)
	if err != nil {
		log.Fatal(err)
	}

	err = m.db.SetPalette(palette)
	if err != nil {
		log.Fatalf("error setting palette: %v", err)
	}
}
