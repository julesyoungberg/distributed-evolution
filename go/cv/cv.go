package cv

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// gets a grayscale image representing the edges of the source image
func GetEdges(img image.Image) (image.Image, error) {
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

// gets the color palette of an image with opencv k means clustering
func GetPalette(img image.Image, k int) ([]color.RGBA, error) {
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
