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
func (m *Master) UpdateTask(args, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if args.Job.ID != m.Job.ID {
		return fmt.Errorf("expected job ID %v, got %v", m.Job.ID, args.Job.ID)
	}

	task := m.Tasks[args.ID]
	if task == nil {
		return fmt.Errorf("task %v not found", args.ID)
	}

	if task.Status == "inprogress" && !task.Connected {
		return fmt.Errorf("disconnected")
	}

	if task.WorkerID == 0 && task.Thread == 0 {
		task.WorkerID = args.WorkerID
		task.Thread = args.Thread
	}

	if args.WorkerID != task.WorkerID || args.Thread != task.Thread {
		return fmt.Errorf("task %v is being worked on by thread %v of worker %v", task.ID, task.Thread, task.WorkerID)
	}

	task.Connected = true
	task.Generation = args.Generation
	task.Status = args.Status
	task.Thread = args.Thread
	task.WorkerID = args.WorkerID
	task.LastUpdate = time.Now()
	m.Tasks[args.ID] = task

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
