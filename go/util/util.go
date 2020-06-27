package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
)

type Image struct {
	Image  image.Image
	Width  int
	Height int
}

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func GetRandomImage() image.Image {
	imageUrl := "https://picsum.photos/900"

	response, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

// EncodeImage encodes an image as png base64 string
func EncodeImage(img image.Image) (string, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		return "", fmt.Errorf("encode: %v", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// DecodeImage decodes a base64 and returns an image
func DecodeImage(data string) (image.Image, error) {
	unbased, _ := base64.StdEncoding.DecodeString(data)
	img, err := png.Decode(bytes.NewReader(unbased))
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}

	return img, nil
}

func GetImageDimensions(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}

func GetSubImage(img image.Image, rect image.Rectangle) image.Image {
	return img.(interface {
		SubImage(rect image.Rectangle) image.Image
	}).SubImage(rect)
}

func RandomColor(rng *rand.Rand) color.RGBA {
	f := func() uint8 {
		return uint8(rng.Intn(4) * 64)
	}

	return color.RGBA{f(), f(), f(), f()}
}

func RandomVector(rng *rand.Rand, bounds Vector) Vector {
	return Vector{X: rng.Float64() * bounds.X, Y: rng.Float64() * bounds.Y}
}

func SquareDifference(x, y float64) float64 {
	d := x - y
	return d * d
}

func RgbaImg(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	return rgba
}

func ImgDiff(ai, bi image.Image) float64 {
	a := RgbaImg(ai)
	b := RgbaImg(bi)
	var d float64
	for i := 0; i < len(a.Pix); i++ {
		d += SquareDifference(float64(a.Pix[i]), float64(b.Pix[i]))
	}

	return math.Sqrt(d)
}
