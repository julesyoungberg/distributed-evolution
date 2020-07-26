package master

import (
	"log"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/gorilla/websocket"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type State struct {
	Fitness          float64                `json:"fitness"`
	Generation       uint                   `json:"generation"`
	JobID            int                    `json:"jobID"`
	NumWorkers       int                    `json:"numWorkers"`
	Output           string                 `json:"output"`
	StartedAt        time.Time              `json:"startedAt"`
	TargetImage      string                 `json:"targetImage"`
	Tasks            map[int]*api.TaskState `json:"tasks"`
	ThreadsPerWorker int                    `json:"threadsPerWorker"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// subscribe handler creates a websocket connection with the client
// TODO handle multiple connections
func (m *Master) subscribe(w http.ResponseWriter, r *http.Request) {
	// TODO check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error establishing websocket connection: ", err)
		return
	}

	log.Printf("new websocket connection request")

	m.connMu.Lock()

	// save the connection for updates
	m.conn = conn

	response := State{TargetImage: m.TargetImageBase64}
	if err := conn.WriteJSON(response); err != nil {
		log.Print("error sending initial state over websocket connection: ", err)
	}

	m.connMu.Unlock()
}

func (m *Master) sendOutput(output *gg.Context, generation uint, fitness float64) {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.conn == nil {
		// no open connections
		return
	}

	m.mu.Lock()
	m.lastUpdate = time.Now()
	tasks := m.Tasks
	jobID := m.Job.ID
	m.mu.Unlock()

	// get resulting image
	img, err := util.EncodeImage(output.Image())
	if err != nil {
		log.Print("[combiner] error sending output: ", err)
		return
	}

	// send encoded image and current generation
	state := State{
		Fitness:    fitness,
		Generation: generation,
		JobID:      jobID,
		Output:     img,
		StartedAt:  m.Job.StartedAt,
		Tasks:      tasks,
	}

	if err := m.conn.WriteJSON(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}
