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
// TODO handle mismatching attempt numbers to handle incorrect duplicate task execution
func (m *Master) Update(args, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task := m.Tasks[args.ID]

	if task.Status == "inprogress" && !task.Connected {
		return fmt.Errorf("disconnected")
	}

	reply.Job.ID = m.Job.ID

	if args.Job.ID != m.Job.ID {
		log.Printf("worker %v is out of date", args.WorkerID)
		return nil
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
