package master

import (
	"fmt"
	"image"
	"log"

	"gocv.io/x/gocv"
)

func getEdges(img image.Image) (image.Image, error) {
	log.Printf("[task-generator] getting edges from image")

	// convert input image to opencv matrix
	src, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return nil, fmt.Errorf("getting edges: %v", err)
	}
	defer src.Close()

	srcGray := gocv.NewMat()
	detectedEdges := gocv.NewMat()
	dst := gocv.NewMatWithSize(src.Rows(), src.Cols(), src.Type())

	lowThreshold := 10
	ratio := 20
	kernelSize := 4

	gocv.CvtColor(src, &srcGray, gocv.ColorBGRToGray)
	gocv.Blur(srcGray, &detectedEdges, image.Point{kernelSize, kernelSize})
	gocv.Canny(detectedEdges, &dst, float32(lowThreshold), float32(lowThreshold*ratio))

	out, err := dst.ToImage()
	if err != nil {
		return nil, fmt.Errorf("getting edges: %v", err)
	}

	return out, err
}
