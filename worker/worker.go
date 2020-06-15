package main

import (
	"time"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Worker struct {
	// worker state
}

func (w *Worker) runTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)
	time.Sleep(10 * time.Second)
}

func main() {
	var w Worker

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for {
		task := getTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.runTask(task)
		}

		time.Sleep(time.Second)
	}
}
