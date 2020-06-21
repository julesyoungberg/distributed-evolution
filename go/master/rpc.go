package master

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
	"github.com/rickyfitts/distributed-evolution/go/worker"
)

func (m *Master) GetTask(args *api.GetTaskArgs, reply *api.Task) error {
	util.DPrintf("task requested")

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

func (m *Master) Update(task, reply *api.Task) error {
	util.DPrintf("update for generation %v received from task %v", task.Generation, task.ID)

	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Job.ID = m.Job.ID

	if task.Job.ID != m.Job.ID {
		log.Printf("worker %v is out of date", task.WorkerID)
		return nil
	}

	task.LastUpdate = time.Now()

	if m.Job.OutputMode == "combined" {
		generation := m.updateGeneration(task)

		m.drawToGeneration(generation, task)

		if generation.Done {
			m.updateUICombined(generation)
			delete(m.Generations, generation.Generation)
		}
	} else {
		overDraw := m.OverDraw

		m.mu.Unlock()

		dc := gg.NewContext(int(task.Dimensions.X)+overDraw*2, int(task.Dimensions.Y)+overDraw*2)
		s := task.BestFit.Genome.(worker.Shapes)
		s.Draw(dc, util.Vector{X: float64(overDraw), Y: float64(overDraw)})
		img := dc.Image()

		m.mu.Lock()

		if o, ok := m.Outputs[task.ID]; !ok || task.Generation > o.Generation {
			m.Outputs[task.ID] = Generation{
				Generation: task.Generation,
				Image:      img,
			}
		}

		m.updateUILatest()
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
