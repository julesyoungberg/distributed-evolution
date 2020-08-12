package master

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
)

// handles a progress update from a worker, updates the state, and updates the ui
func (m *Master) UpdateTask(args, reply *api.TaskState) error {
	m.mu.Lock()

	jobID := m.Job.ID

	t, ok := m.Tasks[args.ID]
	if !ok || t == nil {
		m.mu.Unlock()
		return fmt.Errorf("task %v not found", args.ID)
	}
	task := *t

	nGenerations := m.Job.NumGenerations

	m.mu.Unlock()

	if args.JobID != jobID {
		return fmt.Errorf("expected job ID %v, got %v", jobID, args.JobID)
	}

	if task.Complete {
		return fmt.Errorf("task %v is complete", args.ID)
	}

	if !(task.WorkerID == 0 && task.Thread == 0) && (args.WorkerID != task.WorkerID || args.Thread != task.Thread) {
		return fmt.Errorf("task %v is being worked on by thread %v of worker %v", task.ID, task.Thread, task.WorkerID)
	}

	newTask := *args
	newTask.LastUpdate = time.Now()
	newTask.Attempt = task.Attempt
	newTask.StartedAt = task.StartedAt

	if args.Status == "done" || task.Generation >= nGenerations {
		newTask.Status = "done"
		newTask.Complete = true
		newTask.CompletedAt = time.Now()
	}

	m.mu.Lock()
	m.Tasks[args.ID] = &newTask
	m.mu.Unlock()

	if newTask.Status == "done" && m.allDone() {
		m.mu.Lock()
		m.Job.Complete = true
		m.Job.CompletedAt = time.Now()
		m.mu.Unlock()
	}

	return nil
}

func (m *Master) rpcServer() {
	err := rpc.Register(m)
	if err != nil {
		log.Fatal("rpc error: ", err)
	}

	rpc.HandleHTTP()

	port := os.Getenv("RPC_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("listener error: ", err)
	}

	log.Printf("listening for RPC on port %v\n", port)

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("rpc serve error: ", err)
	}
}
