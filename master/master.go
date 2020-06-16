package main

import (
	"image"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Master struct {
	targetImage       image.Image
	targetImageBase64 string
	taskQueue         []api.Task
	inProgressTasks   []api.Task
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		log.Fatal("error parsing NUM_WORKERS: ", err)
	}

	s := math.Floor(math.Sqrt(float64(numWorkers)))

	width, height := util.GetImageDimensions(m.targetImage)

	cols := int(math.Ceil(float64(width) / s))
	rows := int(math.Ceil(float64(height) / s))

	rgbImg := m.targetImage.(*image.YCbCr)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			x0 := x * int(s)
			y0 := y * int(s)

			x1 := int(math.Min(float64(x0)+s, float64(width)))
			y1 := int(math.Min(float64(y0)+s, float64(height)))

			rect := image.Rect(x0, y0, x1, y1)

			task := api.Task{
				Generation:  1,
				ID:          (y * cols) + x,
				Location:    rect,
				TargetImage: util.EncodeImage(rgbImg.SubImage(rect)),
			}

			m.taskQueue = append(m.taskQueue, task)
		}
	}
}

func main() {
	m := new(Master)

	m.targetImage = util.GetRandomImage()
	m.targetImageBase64 = util.EncodeImage(m.targetImage)
	m.generateTasks()

	m.generateTasks()

	go m.httpServer()

	m.rpcServer()
}
