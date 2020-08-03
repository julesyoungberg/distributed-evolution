package master

import (
	"image"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Result struct {
	ID         int
	Fitness    float64
	Generation uint
	Output     image.Image
	Position   util.Vector
}

// retrieves and decodes the output of a task from the database
func (m *Master) getTaskResult(id int) Result {
	task, err := m.db.GetTask(id)
	if err != nil || task.ID == 0 {
		util.DPrintf("[combiner] error getting task %v", id)
		return Result{ID: -1}
	}

	img, err := util.DecodeImage(task.Output)
	if err != nil {
		util.DPrintf("[combiner] error decoding task %v output", id)
		return Result{ID: -1}
	}

	return Result{
		ID:         task.ID,
		Fitness:    task.BestFit.Fitness,
		Generation: task.Generation,
		Output:     img,
		Position:   task.Position,
	}
}

// listens to the results channel, combining results and sends the final result to the ui
func (m *Master) combineResults(ids []int, results chan Result) {
	total := 0
	var latest uint = 0
	var fitness float64 = 0.0
	dc := gg.NewContext(m.TargetImage.Width, m.TargetImage.Height)

	for result := range results {
		if result.ID < 1 {
			continue
		}

		total++
		fitness += result.Fitness

		if result.Generation > latest {
			latest = result.Generation
		}

		dc.DrawImageAnchored(result.Output, int(result.Position.X), int(result.Position.Y), 0.5, 0.5)
	}

	if total > 0 {
		fitness /= float64(total)
	}

	m.sendOutput(dc, latest, fitness)
}

// periodically read all tasks from db, combine results, save and update ui
func (m *Master) combine() {
	m.mu.Lock()
	jobId := m.Job.ID
	m.mu.Unlock()

	for {
		time.Sleep(5 * time.Second)

		m.mu.Lock()

		if m.Job.ID != jobId {
			jobId = m.Job.ID
			m.mu.Unlock()
			continue
		}

		// get task ids
		ids := []int{}
		for id := range m.Tasks {
			ids = append(ids, id)
		}

		m.mu.Unlock()

		if len(ids) == 0 {
			continue
		}

		var wg sync.WaitGroup
		results := make(chan Result, len(ids))

		wg.Add(len(ids))

		for _, id := range ids {
			go func(id int) {
				results <- m.getTaskResult(id)
				wg.Done()
			}(id)
		}

		go m.combineResults(ids, results)

		wg.Wait()
		close(results)
	}
}
