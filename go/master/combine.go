package master

import (
	"image"
	"log"
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

	if len(task.Output) == 0 {
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
	var generation int64 = 0
	var fitness float64 = 0.0
	dc := gg.NewContext(m.TargetImage.Width, m.TargetImage.Height)

	for result := range results {
		if result.ID < 1 {
			continue
		}

		total++
		generation += int64(result.Generation)
		fitness += result.Fitness

		dc.DrawImageAnchored(result.Output, int(result.Position.X), int(result.Position.Y), 0.5, 0.5)
	}

	if total > 0 {
		generation /= int64(total)
		fitness /= float64(total)
		fitness = 1 / fitness
	}

	m.mu.Lock()
	m.Generation = uint(generation)
	m.Fitness = fitness
	m.mu.Unlock()

	m.sendOutput(dc.Image())
}

// periodically read all tasks from db, combine results, save and update ui
func (m *Master) combine() {
	for {
		time.Sleep(5 * time.Second)

		m.mu.Lock()

		if m.transitioning {
			log.Print("transitioning or complete, sending update")
			m.mu.Unlock()
			m.sendUpdate()
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
