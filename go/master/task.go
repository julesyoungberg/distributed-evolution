package master

import (
	"image"
	"log"
	"math"

	"github.com/google/uuid"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) setTargetImage(image image.Image) {
	width, height := util.GetImageDimensions(image)
	m.TargetImage = util.Image{
		Image:  image,
		Height: height,
		Width:  width,
	}
}

func (m *Master) getTaskRect(x, y, colWidth, rowWidth int) (image.Rectangle, util.Vector) {
	x0 := x * colWidth
	y0 := y * rowWidth
	x1 := int(math.Min(float64(x0+colWidth), float64(m.TargetImage.Width)))
	y1 := int(math.Min(float64(y0+rowWidth), float64(m.TargetImage.Height)))
	return image.Rect(x0, y0, x1, y1), util.Vector{X: float64(x0), Y: float64(y0)}
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	log.Printf("%v workers with %v threads each available, generating tasks...", m.NumWorkers, m.ThreadsPerWorker)

	totalThreads := m.NumWorkers * m.ThreadsPerWorker

	N := math.Floor(math.Sqrt(float64(totalThreads)))
	M := math.Floor(float64(totalThreads) / N)

	colWidth := int(math.Ceil(float64(m.TargetImage.Width) / N))
	rowWidth := int(math.Ceil(float64(m.TargetImage.Height) / M))

	log.Printf("splitting image into %v %vpx cols and %v %vpx rows", N, colWidth, M, rowWidth)

	// create a task for each slice of the image
	for y := 0; y < int(N); y++ {
		for x := 0; x < int(M); x++ {
			rect, offset := m.getTaskRect(x, y, colWidth, rowWidth)

			subImg := util.GetSubImage(m.TargetImage.Image, rect)
			bounds := subImg.Bounds()

			task := api.Task{
				Dimensions:  util.Vector{X: float64(bounds.Dx()), Y: float64(bounds.Dy())},
				Generation:  1,
				ID:          uuid.New().ID(),
				Job:         m.Job,
				Offset:      offset,
				Status:      "unstarted",
				TargetImage: util.EncodeImage(subImg),
				Type:        "polygons",
			}

			err := m.db.PushTask(task)
			if err != nil {
				// shit - TODO try again?
				log.Fatalf("error pushing task to task queue: %v", err)
			}

			m.Tasks[task.ID] = &task
		}
	}

	log.Printf("%v tasks created", len(m.Tasks))
}
