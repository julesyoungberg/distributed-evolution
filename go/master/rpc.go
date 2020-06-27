package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
)

// handles a progress update from a worker, updates the state, and updates the ui
func (m *Master) Update(task, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Job.ID = m.Job.ID

	if task.Job.ID != m.Job.ID {
		log.Printf("worker %v is out of date", task.WorkerID)
		return nil
	}

	task.LastUpdate = time.Now()
	m.Tasks[task.ID] = task

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
