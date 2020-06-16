package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func (m *Master) GetTask(args *api.EmptyMessage, reply *api.Task) error {
	if len(m.taskQueue) > 0 {
		task := m.taskQueue[0]
		task.Started = time.Now()
		*reply = task

		m.inProgressTasks = append(m.inProgressTasks, task)
		m.taskQueue = m.taskQueue[1:]
		util.DPrintf("assigning task %v\n", task.ID)
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
