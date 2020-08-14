package master

import (
	"image/color"
	"log"
	"math/rand"

	"github.com/rickyfitts/distributed-evolution/go/cv"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) savePalete(palette []color.RGBA) {
	log.Printf("[task-generator] saving pallete")

	img := cv.GetPaletteImage(palette)

	encoded, err := util.EncodeImage(img)
	if err != nil {
		log.Printf("[task-generator] failed to encoded palette: %v", err)
		return
	}

	m.mu.Lock()
	m.Palette = encoded
	m.mu.Unlock()

	go m.sendPalette()
}

func (m *Master) getPaletteFromTargetImage(randomCenters bool) []color.RGBA {
	log.Print("[task-generator] getting palette from target image")

	m.mu.Lock()
	img := m.TargetImage.Image
	nColors := m.Job.NumColors
	m.mu.Unlock()

	palette, err := cv.GetPalette(img, nColors, randomCenters)
	if err != nil {
		log.Fatal(err)
	}

	return palette
}

func (m *Master) getRandomPalette() []color.RGBA {
	log.Print("[task-generator] generating random palette")

	palette := make([]color.RGBA, m.Job.NumColors)

	src := rand.NewSource(rand.Int63())
	rng := rand.New(src)

	for i := range palette {
		palette[i] = util.RandomColor(rng)
	}

	return palette
}

func (m *Master) getPalette() []color.RGBA {
	m.mu.Lock()
	paletteType := m.Job.PaletteType
	m.mu.Unlock()

	switch paletteType {
	case "kmeans":
		return m.getPaletteFromTargetImage(false)
	case "kmeansRandomCenters":
		return m.getPaletteFromTargetImage(true)
	case "targetImage":
		return []color.RGBA{{0, 0, 0, 0}} // dummy
	default:
		return m.getRandomPalette()
	}
}

func (m *Master) preparePalette() {
	log.Print("[task-generator] preparing palette")

	palette := m.getPalette()

	err := m.db.SetPalette(palette)
	if err != nil {
		log.Fatalf("error setting palette: %v", err)
	}

	m.savePalete(palette)
}
