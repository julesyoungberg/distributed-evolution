package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
	"net/http"
	"strings"
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
const Debug = 0

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
func EncodeImage(img image.Image) string {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	if err != nil {
		log.Fatal("encode error ", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

// DecodeImage decodes a base64 and returns an image
func DecodeImage(data string) image.Image {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("decode error: ", err)
	}

	return img
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
