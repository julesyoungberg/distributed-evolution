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
func (m *Master) Update(args, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Job.ID = m.Job.ID

	if args.Job.ID != m.Job.ID {
		log.Printf("error! worker %v is out of date", args.WorkerID)
		log.Printf("args.Job.ID: %v, m.Job.ID: %v", args.Job.ID, m.Job.ID)
		return nil
	}

	task := m.Tasks[args.ID]
	if task == nil {
		return nil
	}

	if task.Status == "inprogress" && !task.Connected {
		return fmt.Errorf("disconnected")
	}

	if task.WorkerID == 0 {
		task.WorkerID = args.WorkerID
		reply.WorkerID = args.WorkerID
	} else {
		reply.WorkerID = task.WorkerID

		if args.WorkerID != task.WorkerID {
			log.Printf("error! worker %v claims to be working on task %v, but worker %v is working on it", args.WorkerID, task.ID, task.WorkerID)
			return nil
		}
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
