package master

import (
	"image"
	"log"
	"math"
	"sync"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) allStale() bool {
	for _, t := range m.Tasks {
		if t.Status != "stale" {
			return false
		}
	}

	return true
}

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
	err := m.db.Flush()
	for err != nil {
		log.Printf("[task generator] failed to flush db: %v", err)
		time.Sleep(1 * time.Second)
		err = m.db.Flush()
	}

	m.mu.Lock()

	log.Printf("[task generator] %v workers with %v threads each available", m.NumWorkers, m.ThreadsPerWorker)
	log.Printf("[task generator] generating tasks for job %v", m.Job.ID)

	m.Tasks = map[int]*api.TaskState{}

	totalThreads := m.NumWorkers * m.ThreadsPerWorker

	N := math.Floor(math.Sqrt(float64(totalThreads)))
	M := math.Floor(float64(totalThreads) / N)

	colWidth := int(math.Ceil(float64(m.TargetImage.Width) / N))
	rowWidth := int(math.Ceil(float64(m.TargetImage.Height) / M))

	log.Printf("[task generator] splitting image into %v %vpx cols and %v %vpx rows", N, colWidth, M, rowWidth)

	targetImage := m.TargetImage.Image
	job := m.Job

	m.mu.Unlock()

	var wg sync.WaitGroup

	// create a task for each slice of the image
	for y := 0; y < int(N); y++ {
		for x := 0; x < int(M); x++ {
			wg.Add(1)

			go func(x, y int) {
				defer wg.Done()

				rect, offset := m.getTaskRect(x, y, colWidth, rowWidth)

				subImg := util.GetSubImage(targetImage, rect)
				bounds := subImg.Bounds()

				encoded, err := util.EncodeImage(subImg)
				if err != nil {
					log.Fatal(err)
				}

				task := api.Task{
					Dimensions:  util.Vector{X: float64(bounds.Dx()), Y: float64(bounds.Dy())},
					Generation:  1,
					ID:          y*int(M) + x + 1,
					Job:         job,
					Offset:      offset,
					TargetImage: encoded,
					Status:      "quued",
					Type:        "polygons",
				}

				taskState := api.TaskState{
					ID:         task.ID,
					JobID:      job.ID,
					LastUpdate: time.Now(),
					Status:     "quued",
				}

				m.mu.Lock()
				m.Tasks[task.ID] = &taskState
				m.mu.Unlock()

				e := m.db.PushTask(task)
				if e != nil {
					// let it timeout and try again
					log.Fatalf("[task generator] error pushing task to task queue: %v", e)
					// HACKY - set status to inprogress and let it timeout
					m.Tasks[task.ID].Status = "inprogress"
				}
			}(x, y)
		}
	}

	wg.Wait()

	m.mu.Lock()
	log.Printf("[task generator] %v tasks created", len(m.Tasks))
	m.mu.Unlock()
}

func (m *Master) startRandomTask() {
	log.Print("fetching random image...")
	image := util.GetRandomImage()

	log.Print("encoding image...")
	encodedImg, err := util.EncodeImage(image)
	if err != nil {
		log.Fatal(err)
	}

	m.TargetImageBase64 = encodedImg
	m.setTargetImage(image)

	go m.generateTasks()
}
