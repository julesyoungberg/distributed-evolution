package worker

import (
	"encoding/json"
	"fmt"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) SaveTaskSnapshot() {
	task := w.Task
	task.Population = w.ga.Populations[0]

	encoded, err := json.Marshal(task)
	if err != nil {
		fmt.Print("error encoding task snapshot: ", err)
		return
	}

	err = w.cache.Set(util.GetSnapshotKey(task.ID), string(encoded))
	if err != nil {
		fmt.Print("error saving snapshot: ", err)
	}
}
