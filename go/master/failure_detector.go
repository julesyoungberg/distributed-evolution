package master

import (
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// recovers a failed task by fetching the
func (m *Master) recoverTask(id int) {
	m.mu.Lock()

	task := m.Tasks[id]

	m.mu.Unlock()

	snapshot, err := m.GetTaskSnapshot(task.ID)
	if err != nil {
		log.Fatalf("error recovering task %v: %v", task.ID, err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Job.ID != snapshot.Job.ID {
		log.Fatalf("unable to recover task %v: job ID is out of date")
	}

	util.DPrintf("fetched worker snapshot for task %v", task.ID)
	task = snapshot
	task.Status = "recovering"

	var fastest *api.Task = nil

	for _, t := range m.Tasks {
		if t.ID == task.ID {
			continue
		}

		if fastest == nil || t.Generation > fastest.Generation {
			*fastest = t
		}
	}

	fastest.Linked = append(fastest.Linked, task.Linked...)
	task.Linked = fastest.Linked

	m.Tasks[task.ID] = task
	m.Tasks[fastest.ID] = *fastest
}

// detect worker failures
func (m *Master) DetectFailures() {
	for {
		time.Sleep(time.Second)

		m.mu.Lock()

		for i, t := range m.Tasks {
			if time.Since(t.LastUpdate) > time.Second {
				util.DPrintf("task %v timed out! recovering...", i)
				m.Tasks[i].Status = "failed"
				go m.recoverTask(i)
			}
		}

		m.mu.Unlock()
	}
}
