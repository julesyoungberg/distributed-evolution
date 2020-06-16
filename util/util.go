package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strings"
)

// Debugging
const Debug = 1

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

func GetRandomImage() image.Image {
	DPrintf("fetching random image...")

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

func EncodeImage(img image.Image) string {
	buf := new(bytes.Buffer)

	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func DecodeImage(data string) image.Image {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("error decoding task target image ", err)
	}

	return img
}

func GetImageDimensions(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy()
}
