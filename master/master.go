package main

import (
	"image"
	"math"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Master struct {
	targetImage       image.Image
	targetImageBase64 string
	taskQueue         []api.Task
	inProgressTasks   []api.Task
	finishedTasks     []api.Task
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	sliceSize := 100
	bounds := m.targetImage.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	cols := int(math.Ceil(float64(width) / float64(sliceSize)))
	rows := int(math.Ceil(float64(height) / float64(sliceSize)))

	rgbImg := m.targetImage.(*image.YCbCr)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			x0 := x * sliceSize
			y0 := y * sliceSize

			x1 := int(math.Min(float64(x0+sliceSize), float64(width)))
			y1 := int(math.Min(float64(y0+sliceSize), float64(height)))

			rect := image.Rect(x0, y0, x1, y1)

			task := api.Task{
				Generation:  1,
				ID:          (y * cols) + x,
				Location:    rect,
				TargetImage: util.Base64EncodeImage(rgbImg.SubImage(rect)),
			}

			m.taskQueue = append(m.taskQueue, task)
		}
	}
}

func main() {
	m := new(Master)

	m.targetImage = util.GetRandomImage()
	m.targetImageBase64 = util.Base64EncodeImage(m.targetImage)
	m.generateTasks()

	m.generateTasks()

	go m.httpServer()

	m.rpcServer()
}
