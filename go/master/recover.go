package master

import (
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

// recovers a failed task by fetching the
func (m *Master) recover(id int) {
	m.mu.Lock()
	m.Tasks[id].Status = "recovering"
	m.Tasks[id].Connected = true
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
// if a task timesout, mark it as failed and begin recovery process
func (m *Master) detectFailures() {
	timeout := 30 * time.Second
	queueTimeout := timeout * 2

	for {
		time.Sleep(timeout / 2)

		m.mu.Lock()

		for i, t := range m.Tasks {
			since := time.Since(t.LastUpdate)
			timedout := t.Status == "inprogress" && since > timeout
			timedout = timedout || (t.Status == "queued" && since > queueTimeout)

			if timedout {
				util.DPrintf("[failure detector] task %v timed out! recovering...", i)
				m.Tasks[i].Status = "failed"
				go m.recover(i)
			}
		}

		m.mu.Unlock()
	}
}
