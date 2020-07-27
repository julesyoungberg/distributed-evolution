package master

import (
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

// recovers a failed task by fetching the
func (m *Master) recover(id int) {
	m.mu.Lock()
	task := *m.Tasks[id]
	jobID := m.Job.ID
	m.mu.Unlock()

	// in case an update was received between
	// this function being called and the lock being received
	if task.Status != "failed" {
		return
	}

	// check if this task is still up to date
	if task.JobID != jobID {
		// if this task is from a previous job,
		// mark it as stale and forget about it
		log.Printf("[failure detector] task %v is from job %v, the current job is %v", task.ID, task.JobID, m.Job.ID)

		m.mu.Lock()
		m.Tasks[id].Status = "stale"
		m.mu.Unlock()

		return
	}

	task.Status = "queued"
	task.WorkerID = 0
	task.Thread = 0
	task.LastUpdate = time.Now()

	m.mu.Lock()
	m.Tasks[id] = &task
	m.mu.Unlock()

	err := m.db.PushTaskID(task.ID)
	if err != nil {
		log.Printf("[failure detector] error pushing task: %v", err)

		// try again next round
		m.mu.Lock()
		task.Status = "inprogress"
		m.mu.Unlock()
	} else {
		log.Printf("[failure detector] requeued task %v", task.ID)
	}
}

// check that each inprogress task is active by checking its last update
// if a task times out, mark it as failed and begin recovery process
func (m *Master) detectFailures() {
	timeout := 20 * time.Second
	queueTimeout := 60 * time.Second

	for {
		time.Sleep(timeout / 4)

		m.mu.Lock()

		for id, task := range m.Tasks {
			workerTimeout := task.Status == "inprogress" && time.Since(task.LastUpdate) > timeout
			queueTimeout := task.Status == "queued" && time.Since(task.LastUpdate) > queueTimeout

			if workerTimeout || queueTimeout {
				util.DPrintf("[failure detector] %v task %v timed out, recovering", task.Status, id)
				task.Status = "failed"
				go m.recover(id)
			}
		}

		m.mu.Unlock()
	}
}
