package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// assigns a task to a worker
func (m *Master) GetTask(args *api.GetTaskArgs, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// find an unstarted task to reply with
	for i, task := range m.Tasks {
		if task.Status == "unstarted" {
			task.Status = "active"
			task.WorkerID = args.WorkerID
			task.LastUpdate = time.Now()

			m.Tasks[i] = task
			*reply = task
			reply.Job = m.Job

			log.Printf("assigning task %v job %v to worker %v\n", task.ID, m.Job.ID, args.WorkerID)

			return nil
		}
	}

	return nil
}

func (m *Master) addNewLinkedTasks(task, reply *api.Task) {
	newTasks := []api.Task{}
	linked := m.Tasks[task.ID].Linked

	if len(task.Linked) < len(linked) {
		for _, i := range linked {
			t := m.Tasks[i]

			if t.Status == "recovering" {
				util.DPrintf("assigning task %v to worker %v", t.ID, task.WorkerID)
				t.Status = "active"
				t.WorkerID = task.WorkerID
				t.LastUpdate = time.Now()

				newTasks = append(newTasks, t)
				m.Tasks[i] = t
			}
		}

		task.Linked = linked
		reply.Linked = linked
	}
}

// handles a progress update from a worker, updates the state, and updates the ui
func (m *Master) Update(task, reply *api.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Job.ID = m.Job.ID

	if task.Job.ID != m.Job.ID {
		log.Printf("worker %v is out of date", task.WorkerID)
		return nil
	}

	m.addNewLinkedTasks(task, reply)

	task.LastUpdate = time.Now()

	if m.Job.OutputMode == "combined" {
		// draw the output to the corresponding generation
		generation := m.updateGeneration(task)
		m.drawToGeneration(&generation, task)
		m.Generations[generation.Generation] = generation

		if generation.Done {
			m.updateUICombined(generation)
			delete(m.Generations, generation.Generation)
		}
	} else {
		if time.Since(m.lastUpdate) > m.wsHeartbeatTimeout/2.0 {
			m.updateUILatest()
		}
	}

	m.Tasks[task.ID] = *task

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
