package master

import (
	"image"
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

// assigns a task to a worker
func (m *Master) GetTask(args *api.GetTaskArgs, reply *api.Task) error {
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

	if m.Job.OutputMode == "combined" {
		// draw the output to the corresponding generation
		generation := m.updateGeneration(task)

		m.drawToGeneration(generation, task)

		if generation.Done {
			m.updateUICombined(generation)
			delete(m.Generations, generation.Generation)
		}
	} else {
		// save the output to the outputs map if the generation is the latest
		if o, ok := m.Outputs[task.ID]; ok && task.Generation < o.Generation {
			// this update is old, must have been delayed, ignore
			return
		}

		var img image.Image

		if m.Job.DrawOnce {
			// the worker has already drawn the generation, use that
			img = util.DecodeImage(task.Output)
		} else {
			// draw the generation
			overDraw := m.Job.OverDraw

			m.mu.Unlock()

			dc := gg.NewContext(int(task.Dimensions.X)+overDraw*2, int(task.Dimensions.Y)+overDraw*2)
			s := task.BestFit.Genome.(worker.Shapes)
			s.Draw(dc, util.Vector{X: float64(overDraw), Y: float64(overDraw)})
			img = dc.Image()

			m.mu.Lock()
		}

		// save the output
		m.Outputs[task.ID] = Generation{
			Generation: task.Generation,
			Image:      img,
		}

		// TODO throttle this to increase performance
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
