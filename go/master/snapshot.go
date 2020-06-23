package master

import (
	"encoding/json"
	"fmt"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// fetches and parses a task snapshot from the cache
func (m *Master) GetTaskSnapshot(id int) (api.Task, error) {
	val, err := m.cache.Get(util.GetSnapshotKey(id))
	if err != nil {
		e := fmt.Errorf("error fetching snapshot for worker %v: %v", id, err)
		return api.Task{}, e
	}

	bytes := []byte(val)

	var snapshot api.Task

	err = json.Unmarshal(bytes, &snapshot)
	if err != nil {
		e := fmt.Errorf("error parsing snapshot for worker %v: %v", id, err)
		return api.Task{}, e
	}

	return snapshot, nil
}
