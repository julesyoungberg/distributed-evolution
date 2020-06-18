package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func (m *Master) GetTask(args, reply *api.Task) error {
	util.DPrintf("task requested")

	for _, task := range m.Tasks {
		if task.Status == "unstarted" {
			task.Status = "active"
			*reply = task
			util.DPrintf("assigning task %v\n", task.ID)
		}
	}

	return nil
}

func (m *Master) Update(task, reply *api.Task) error {
	util.DPrintf("update received from task %v", task.ID)

	m.mu.Lock()
	m.Tasks[task.ID] = *task
	m.mu.Unlock()

	genN := m.updateGenerations(*task)

	m.drawToGeneration(genN, *task)

	m.updateUI(genN)

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
