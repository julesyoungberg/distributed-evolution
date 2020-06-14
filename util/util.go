package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
)

func GetRandomImage() image.Image {
	fmt.Printf("fetching image...")

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
	fmt.Printf("encoding image...")

	buf := new(bytes.Buffer)

	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
