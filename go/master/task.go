package master

import (
	"image"
	"log"
	"math"
	"sync"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/cv"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) getTaskRect(x, y, colWidth, rowWidth int) (image.Rectangle, util.Vector) {
	m.mu.Lock()
	targetImage := m.TargetImage
	overDraw := m.Job.OverDraw
	m.mu.Unlock()

	x0 := x * colWidth
	if x0 > 0 {
		x0 -= overDraw
	}

	y0 := y * rowWidth
	if y0 > 0 {
		y0 -= overDraw
	}

	x1 := int(math.Min(float64(x0+colWidth), float64(targetImage.Width)))
	if x1 < targetImage.Width {
		x1 += overDraw
	}

	y1 := int(math.Min(float64(y0+rowWidth), float64(targetImage.Height)))
	if y1 < targetImage.Height {
		y1 += overDraw
	}

	center := util.Vector{X: float64(x0 + (x1-x0)/2), Y: float64(y0 + (y1-y0)/2)}

	return image.Rect(x0, y0, x1, y1), center
}

// generate a task given configuration, save it in local state and push to task queue
func (m *Master) generateTask(
	x, y, colWidth, rowWidth, M int,
	targetImage image.Image,
	edges image.Image,
	job api.Job,
) {
	rect, pos := m.getTaskRect(x, y, colWidth, rowWidth)

	subImg := util.GetSubImage(targetImage, rect)
	encoded, err := util.EncodeImage(subImg)
	if err != nil {
		log.Fatal(err)
	}

	task := api.Task{
		Dimensions:         util.Vector{X: float64(colWidth), Y: float64(rowWidth)},
		Generation:         1,
		ID:                 y*M + x + 1,
		Job:                job,
		Position:           pos,
		TargetImage:        encoded,
		ScaledQuantization: job.Quantization / M,
		ShapeType:          job.ShapeType,
	}

	if edges != nil {
		subImgEdges := util.GetSubImage(edges, rect)
		encodedEdges, err := util.EncodeImage(subImgEdges)
		if err != nil {
			log.Fatal(err)
		}

		task.Edges = encodedEdges
	}

	taskState := api.TaskState{
		Attempt:    1,
		ID:         task.ID,
		JobID:      job.ID,
		LastUpdate: time.Now(),
		StartedAt:  time.Now(),
		Status:     "queued",
	}

	m.mu.Lock()
	m.Tasks[task.ID] = &taskState
	m.mu.Unlock()

	e := m.db.PushTask(task)
	if e != nil {
		// let it timeout and try again
		log.Fatalf("[task-generator] error pushing task to task queue: %v", e)
		m.mu.Lock()
		// HACKY - set status to inprogress and let it timeout
		m.Tasks[task.ID].Status = "inprogress"
		m.mu.Unlock()
	}
}

func (m *Master) saveEdges(edges image.Image) {
	log.Print("[task-generator] saving edges")

	encodedEdges, err := util.EncodeImage(edges)
	if err != nil {
		log.Printf("[task-generator] failed to encode edges: %v", err)
	}

	m.mu.Lock()
	m.TargetImageEdges = encodedEdges
	m.mu.Unlock()

	m.sendEdges()
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	log.Print("[task-generator] flushing the db")

	err := m.db.Flush()
	for err != nil {
		log.Printf("[task-generator] failed to flush db: %v", err)
		time.Sleep(1 * time.Second)
		err = m.db.Flush()
	}

	m.mu.Lock()

	log.Printf("[task-generator] %v workers with %v threads each available", m.NumWorkers, m.ThreadsPerWorker)
	log.Printf("[task-generator] generating tasks for job %v", m.Job.ID)

	m.Tasks = map[int]*api.TaskState{}

	totalThreads := m.NumWorkers * m.ThreadsPerWorker

	N := math.Floor(math.Sqrt(float64(totalThreads)))
	M := math.Floor(float64(totalThreads) / N)
	m.Job.ShapesPerSlice = m.Job.NumShapes / (int(N) * int(M))

	colWidth := int(math.Ceil(float64(m.TargetImage.Width) / N))
	rowWidth := int(math.Ceil(float64(m.TargetImage.Height) / M))

	log.Printf("[task-generator] splitting image into %v %vpx cols and %v %vpx rows", N, colWidth, M, rowWidth)

	targetImage := m.TargetImage.Image
	job := m.Job

	m.mu.Unlock()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		m.preparePalette()
		wg.Done()
	}()

	var edges image.Image

	if job.DetectEdges {
		log.Printf("[task-generator] getting target image edges")
		edges, err = cv.GetEdges(targetImage)
		if err != nil {
			log.Fatal(err)
		}

		go m.saveEdges(edges)
	}

	// create a task for each slice of the image
	for y := 0; y < int(N); y++ {
		for x := 0; x < int(M); x++ {
			wg.Add(1)

			go func(x, y int) {
				m.generateTask(x, y, colWidth, rowWidth, int(M), targetImage, edges, job)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()

	m.mu.Lock()
	nTasks := len(m.Tasks)
	m.mu.Unlock()

	log.Printf("[task-generator] %v tasks created", nTasks)
}
