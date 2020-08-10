package cv

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
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

	lowThreshold := 8
	ratio := 20
	kernelSize := 8

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
func GetPalette(img image.Image, k int, randomCenters bool) ([]color.RGBA, error) {
	// convert input image to opencv matrix
	src, err := gocv.ImageToMatRGB(img)
	if err != nil {
		return []color.RGBA{}, fmt.Errorf("getting palette: %v", err)
	}
	defer src.Close()

	data := src.Reshape(1, src.Total())
	data.ConvertTo(&data, gocv.MatTypeCV32F)

	labels := gocv.NewMat()
	criteria := gocv.NewTermCriteria(gocv.MaxIter, 20, 1.0)
	centers := gocv.NewMat()

	// this appears to a bug in gocv,
	// where these two flags are oppositely valued compared to opencv
	flags := gocv.KMeansUseInitialLabels // gocv.KMeansPPCenters
	if randomCenters {
		flags = gocv.KMeansRandomCenters
	}

	log.Printf("calling kmeans with k: %v, flags: %v (%v)", k, flags, int(flags))

	gocv.KMeans(data, k, &labels, criteria, 1, flags, &centers)

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

func GetPaletteImage(palette []color.RGBA) image.Image {
	nColors := len(palette)
	colorsPerEdge := int(math.Ceil(math.Sqrt(float64(nColors))))
	size := 32

	dc := gg.NewContext(size*colorsPerEdge, size*colorsPerEdge)

	for i, color := range palette {
		y := i / colorsPerEdge
		x := i % colorsPerEdge

		dc.DrawRectangle(float64(x*size), float64(y*size), float64(size), float64(size))
		dc.SetColor(color)
		dc.Fill()
	}

	return dc.Image()
}
