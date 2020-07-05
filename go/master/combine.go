package master

import (
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

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

		dc := gg.NewContext(m.TargetImage.Width, m.TargetImage.Height)

		var latest uint = 0
		var fitness float64 = 0.0

		// combine outputs
		for _, id := range ids {
			task, err := m.db.GetTask(id)
			if err != nil || task.ID == 0 {
				continue
			}

			if task.Generation > latest {
				latest = task.Generation
			}

			fitness += task.BestFit.Fitness

			img, err := util.DecodeImage(task.Output)
			if err != nil {
				continue
			}

			centerX := int(math.Round(task.Offset.X + task.Dimensions.X/2.0))
			centerY := int(math.Round(task.Offset.Y + task.Dimensions.Y/2.0))

			dc.DrawImageAnchored(img, centerX, centerY, 0.5, 0.5)
		}

		fitness /= float64(len(ids))

		m.sendOutput(dc, latest, fitness)
	}
}
