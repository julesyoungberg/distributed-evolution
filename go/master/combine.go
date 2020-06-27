package master

import (
	"log"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// periodically read all tasks from db, combine results, save and update ui
func (m *Master) combine() {
	for {
		time.Sleep(5 * time.Second)

		m.mu.Lock()

		dc := gg.NewContext(m.TargetImage.Width, m.TargetImage.Height)

		var latest uint = 0

		log.Printf("[combiner] combining outputs")

		for id := range m.Tasks {
			m.mu.Unlock()

			task, err := m.db.GetTask(id)
			if err != nil {
				continue
			}

			if task.Generation > latest {
				latest = task.Generation
			}

			img, err := util.DecodeImage(task.Output)
			if err != nil {
				log.Print("error: ", err)
				continue
			}

			centerX := int(math.Round(task.Offset.X + task.Dimensions.X/2.0))
			centerY := int(math.Round(task.Offset.Y + task.Dimensions.Y/2.0))

			dc.DrawImageAnchored(img, centerX, centerY, 0.5, 0.5)

			m.mu.Lock()
		}

		m.mu.Unlock()

		log.Printf("[combiner] sending output")

		m.sendOutput(dc, latest)
	}
}
