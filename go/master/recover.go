package master

import (
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

// recovers a failed task by fetching the
func (m *Master) recover(id uint32) {
	m.mu.Lock()
	task := m.Tasks[id]
	task.Status = "recovering"
	m.mu.Unlock()

	err := m.db.PushTaskID(task.ID)
	if err != nil {
		log.Printf("error pushing task to task queue: %v", err)

		m.mu.Lock()
		task.Status = "inprogress"
		m.mu.Unlock()
	}
}

// check that each inprogress task is active by checking its last update
// if a task timesout, mark it as failed and begin recovery process
func (m *Master) detectFailures() {
	timeout := 20 * time.Second

	for {
		time.Sleep(time.Second)

		m.mu.Lock()

		for i, t := range m.Tasks {
			if t.Status == "inprogress" && time.Since(t.LastUpdate) > timeout {
				util.DPrintf("task %v timed out! recovering...", i)
				m.Tasks[i].Status = "failed"
				go m.recover(i)
			}
		}

		m.mu.Unlock()
	}
}
