package worker

import (
	"encoding/json"
	"fmt"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/cache"
)

type WorkerSnapshot struct {
	ID           uint32    `json:"ID"`
	GA           *eaopt.GA `json:"GA"`
	NGenerations uint      `json:"nGenerations"`
	Task         api.Task  `json:"task"`
}

func getSnapshotKey(id uint32) string {
	return fmt.Sprintf("workerID%v", id)
}

// saves a worker snapshot as a serialized JSON string to the cache
func (w *Worker) SaveWorkerSnapshot() {
	snapshot := WorkerSnapshot{
		ID:           w.ID,
		GA:           w.ga,
		NGenerations: w.NGenerations,
		Task:         w.Task,
	}

	encoded, err := json.Marshal(snapshot)
	if err != nil {
		fmt.Print("error encoding snapshot: ", err)
		return
	}

	err = w.cache.Set(getSnapshotKey(w.ID), string(encoded))
	if err != nil {
		fmt.Print("error saving snapshot: ", err)
	}
}

// fetches and parses a worker snapshot from the cache
func GetWorkerSnapshot(cache *cache.Cache, id uint32) (WorkerSnapshot, error) {
	val, err := cache.Get(getSnapshotKey(id))
	if err != nil {
		e := fmt.Errorf("error fetching snapshot for worker %v: %v", id, err)
		return WorkerSnapshot{}, e
	}

	bytes := []byte(val)

	var snapshot WorkerSnapshot

	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		e := fmt.Errorf("error parsing snapshot for worker %v: %v", id, err)
		return WorkerSnapshot{}, e
	}

	return snapshot, nil
}
