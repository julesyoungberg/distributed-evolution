package master

import (
	"image"
	"image/color"
	"log"
	"math"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/cv"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) allStale() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("checking for staleness")

	for _, t := range m.Tasks {
		if t.Status != "stale" {
			log.Printf("not all tasks are stale, task %v is %v", t.ID, t.Status)
			return false
		}
	}

	log.Printf("all tasks are stale")
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
	x, y, colWidth, rowWidth int,
	M float64,
	targetImage image.Image,
	edges image.Image,
	job api.Job,
) {
	rect, pos := m.getTaskRect(x, y, colWidth, rowWidth)

	subImg := util.GetSubImage(targetImage, rect)
	subImgEdges := util.GetSubImage(edges, rect)

	encoded, err := util.EncodeImage(subImg)
	if err != nil {
		log.Fatal(err)
	}

	encodedEdges, err := util.EncodeImage(subImgEdges)
	if err != nil {
		log.Fatal(err)
	}

	task := api.Task{
		Dimensions:  util.Vector{X: float64(colWidth), Y: float64(rowWidth)},
		Edges:       encodedEdges,
		Generation:  1,
		ID:          y*int(M) + x + 1,
		Job:         job,
		Position:    pos,
		TargetImage: encoded,
		ShapeType:   m.Job.ShapeType,
	}

	taskState := api.TaskState{
		ID:         task.ID,
		JobID:      job.ID,
		LastUpdate: time.Now(),
		Status:     "queued",
	}

	m.mu.Lock()
	m.Tasks[task.ID] = &taskState
	m.mu.Unlock()

	e := m.db.PushTask(task)
	if e != nil {
		// let it timeout and try again
		log.Fatalf("[task generator] error pushing task to task queue: %v", e)
		m.mu.Lock()
		// HACKY - set status to inprogress and let it timeout
		m.Tasks[task.ID].Status = "inprogress"
		m.mu.Unlock()
	}
}

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

func (m *Master) preparePalette() {
	log.Printf("[task-generator] preparing palette")

	m.mu.Lock()
	img := m.TargetImage.Image
	nColors := m.Job.NumColors
	m.mu.Unlock()

	palette, err := cv.GetPalette(img, nColors)
	if err != nil {
		log.Fatal(err)
	}

	err = m.db.SetPalette(palette)
	if err != nil {
		log.Fatalf("error setting palette: %v", err)
	}

	m.savePalete(palette)
}

func (m *Master) saveEdges(edges image.Image) {
	log.Print("[task-generator] saving edges")

	encodedEdges, err := util.EncodeImage(edges)
	if err != nil {
		log.Printf("[task generator] failed to encode edges: %v", err)
	}

	m.mu.Lock()
	m.TargetImageEdges = encodedEdges
	m.mu.Unlock()

	m.sendEdges()
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
	m.Job.ShapesPerSlice = m.Job.NumShapes / (int(N) * int(M))

	colWidth := int(math.Ceil(float64(m.TargetImage.Width) / N))
	rowWidth := int(math.Ceil(float64(m.TargetImage.Height) / M))

	log.Printf("[task generator] splitting image into %v %vpx cols and %v %vpx rows", N, colWidth, M, rowWidth)

	targetImage := m.TargetImage.Image
	job := m.Job

	m.mu.Unlock()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		m.preparePalette()
		wg.Done()
	}()

	log.Printf("[task-generator] getting target image edges")
	edges, err := cv.GetEdges(targetImage)
	if err != nil {
		log.Fatal(err)
	}

	go m.saveEdges(edges)

	// create a task for each slice of the image
	for y := 0; y < int(N); y++ {
		for x := 0; x < int(M); x++ {
			wg.Add(1)

			go func(x, y int) {
				m.generateTask(x, y, colWidth, rowWidth, M, targetImage, edges, job)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()

	m.mu.Lock()
	nTasks := len(m.Tasks)
	m.mu.Unlock()

	log.Printf("[task generator] %v tasks created", nTasks)
}

func (m *Master) startRandomJob() {
	log.Print("fetching random image...")
	image := util.GetRandomImage()

	log.Print("encoding image...")
	encodedImg, err := util.EncodeImage(image)
	if err != nil {
		log.Fatal(err)
	}

	m.TargetImageBase64 = encodedImg
	m.setTargetImage(image)
	m.Job.StartedAt = time.Now()

	go m.generateTasks()
}

func (m *Master) stopJob() {
	log.Printf("transitioning")

	m.mu.Lock()
	m.transitioning = true
	m.mu.Unlock()

	log.Printf("clearing the task queue")

	// clear the queue of all tasks
	for {
		_, err := m.db.PullTask()
		if err != nil {
			break
		}
	}

	m.mu.Lock()

	// mark any qeued tasks as stale
	for _, task := range m.Tasks {
		if task.Status == "queued" {
			task.Status = "stale"
		}
	}

	m.mu.Unlock()

	log.Printf("waiting for workers to stop")

	// wait for all the workers to stop
	for !m.allStale() {
		time.Sleep(time.Second)
	}

	m.mu.Lock()
	m.transitioning = false
	m.mu.Unlock()

	log.Printf("transition complete")
}
