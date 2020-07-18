package master

import (
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) saveSnapshots() {
	for {
		time.Sleep(30 * time.Second)

		snapshot := api.MasterSnapshot{
			Job:         m.Job,
			TargetImage: m.TargetImageBase64,
			Tasks:       make(map[int]api.TaskState, len(m.Tasks)),
		}

		for id, task := range m.Tasks {
			snapshot.Tasks[id] = *task
		}

		err := m.db.SaveSnapshot(snapshot)
		if err != nil {
			util.DPrintf("[snapshotter] error: %v", err)
		}
	}
}

func (m *Master) restoreFromSnapshot() bool {
	snapshot, err := m.db.GetSnapshot()
	if err != nil {
		util.DPrintf("error restoring from snapshot: %v", err)
		return false
	}

	image, err := util.DecodeImage(m.TargetImageBase64)
	if err != nil {
		util.DPrintf("error decoding target image from snapshot: %v", err)
		return false
	}

	m.Job = snapshot.Job
	m.TargetImageBase64 = snapshot.TargetImage
	m.setTargetImage(image)
	m.Tasks = make(map[int]*api.TaskState, len(snapshot.Tasks))

	for id, task := range snapshot.Tasks {
		m.Tasks[id] = &task
	}

	return true
}
