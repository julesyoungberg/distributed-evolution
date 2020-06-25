package worker

import (
	"encoding/json"
	"fmt"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) SaveTaskSnapshot(state *WorkerTask) {
	task := state.Task
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

// fetches and parses a task snapshot from the cache
func (w *Worker) GetTaskSnapshot(id int) (api.Task, error) {
	val, err := w.cache.Get(util.GetSnapshotKey(id))
	if err != nil {
		e := fmt.Errorf("error fetching snapshot for task %v: %v", id, err)
		return api.Task{}, e
	}

	bytes := []byte(val)

	var snapshot api.Task

	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		e := fmt.Errorf("error parsing snapshot for task%v: %v", id, err)
		return api.Task{}, e
	}

	return snapshot, nil
}
