package master

import (
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

// recovers a failed task by fetching the
func (m *Master) recover(id int) {
	m.mu.Lock()
	task := m.Tasks[id]

	// in case an update was received between
	// this function being called and the lock being received
	if task.Status != "failed" {
		m.mu.Unlock()
		return
	}

	// check if this task is still up to date
	if task.JobID != m.Job.ID {
		// if this task is from a previous job,
		// mark it as stale and forget about it
		log.Printf("[failure detector] task %v is from job %v, the current job is %v", task.ID, task.JobID, m.Job.ID)
		task.Status = "stale"
		m.Tasks[id] = task
		m.mu.Unlock()
		return
	}

	task.Status = "queued"
	task.WorkerID = 0
	task.Thread = 0
	m.mu.Unlock()

	err := m.db.PushTaskID(m.Tasks[id].ID)
	if err != nil {
		log.Printf("[failure detector] error %v", err)

		// try again next round
		m.mu.Lock()
		m.Tasks[id].Status = "inprogress"
		m.mu.Unlock()
	}
}

// check that each inprogress task is active by checking its last update
// if a task times out, mark it as failed and begin recovery process
func (m *Master) detectFailures() {
	timeout := 5 * time.Second
	queueTimeout := 30 * time.Second

	for {
		time.Sleep(timeout / 4)

		m.mu.Lock()

		for i, t := range m.Tasks {
			workerTimeout := t.Status == "inprogress" && time.Since(t.LastUpdate) > timeout
			queueTimeout := t.Status == "queued" && time.Since(t.LastUpdate) > queueTimeout

			if workerTimeout || queueTimeout {
				util.DPrintf("[failure detector] task %v timed out!", i)
				m.Tasks[i].Status = "failed"
				go m.recover(i)
			}
		}

		m.mu.Unlock()
	}
}
