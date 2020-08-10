package master

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/cv"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) savePalete(palette []color.RGBA) {
	log.Printf("[task-generator] saving pallete")
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

	img := dc.Image()
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

func (m *Master) getPaletteFromTargetImage() []color.RGBA {
	log.Print("[task-generator] getting palette from target image")

	m.mu.Lock()
	img := m.TargetImage.Image
	nColors := m.Job.NumColors
	m.mu.Unlock()

	palette, err := cv.GetPalette(img, nColors)
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
	case "targetImage":
		return m.getPaletteFromTargetImage()
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
