package master

import (
	"image"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type State struct {
	Complete         bool                  `json:"complete"`
	CompletedAt      time.Time             `json:"completedAt"`
	Fitness          float64               `json:"fitness"`
	Generation       uint                  `json:"generation"`
	JobID            int                   `json:"jobID"`
	NumWorkers       int                   `json:"numWorkers"`
	Output           string                `json:"output"`
	Palette          string                `json:"palette"`
	StartedAt        time.Time             `json:"startedAt"`
	TargetImage      string                `json:"targetImage"`
	TargetImageEdges string                `json:"targetImageEdges"`
	Tasks            map[int]api.TaskState `json:"tasks"`
	ThreadsPerWorker int                   `json:"threadsPerWorker"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (m *Master) getState() State {
	m.mu.Lock()

	tasks := make(map[int]api.TaskState, len(m.Tasks))

	for id, task := range m.Tasks {
		tasks[id] = *task
	}

	// send encoded image and current generation
	state := State{
		Complete:         m.Job.Complete,
		CompletedAt:      m.Job.CompletedAt,
		Fitness:          m.Fitness,
		Generation:       m.Generation,
		JobID:            m.Job.ID,
		NumWorkers:       m.NumWorkers,
		StartedAt:        m.Job.StartedAt,
		ThreadsPerWorker: m.ThreadsPerWorker,
		Tasks:            tasks,
	}

	m.mu.Unlock()

	return state
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

	state := m.getState()

	m.mu.Lock()
	state.Palette = m.Palette
	state.TargetImage = m.TargetImageBase64
	state.TargetImageEdges = m.TargetImageEdges
	m.mu.Unlock()

	// save the connection for updates
	if err := conn.WriteJSON(state); err != nil {
		log.Print("error sending initial state over websocket connection: ", err)
		return
	}

	m.connMu.Lock()
	m.conn = conn
	m.connMu.Unlock()
}

func (m *Master) sendData(state State) error {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.conn == nil {
		return nil
	}

	return m.conn.WriteJSON(state)
}

func (m *Master) sendOutput(output image.Image) {
	state := m.getState()

	if img, err := util.EncodeImage(output); err == nil {
		state.Output = img
	} else {
		log.Printf("[combiner] error encoding output image: %v", err)
	}

	if err := m.sendData(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}

func (m *Master) sendUpdate() {
	state := m.getState()

	if err := m.sendData(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}

func (m *Master) sendEdges() {
	log.Printf("[task-generator] sending edges")

	state := m.getState()

	m.mu.Lock()
	state.TargetImageEdges = m.TargetImageEdges
	m.mu.Unlock()

	if err := m.sendData(state); err != nil {
		log.Printf("[task-generator] error sending edges: %v", err)
	}
}

func (m *Master) sendPalette() {
	log.Printf("[task-generator] sending palette")

	state := m.getState()

	m.mu.Lock()
	state.Palette = m.Palette
	m.mu.Unlock()

	if err := m.sendData(state); err != nil {
		log.Printf("[task-generator] error sending palette: %v", err)
	}
}
