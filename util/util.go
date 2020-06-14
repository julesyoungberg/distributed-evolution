package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"log"
	"net/http"
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

func Base64EncodeImage(img image.Image) string {
	DPrintf("encoding image...")

	buf := new(bytes.Buffer)

	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
