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

	m.mu.Lock()

	response := State{
		JobID:            m.Job.ID,
		Palette:          m.Palette,
		TargetImage:      m.TargetImageBase64,
		TargetImageEdges: m.TargetImageEdges,
	}

	m.mu.Unlock()

	// save the connection for updates
	if err := conn.WriteJSON(response); err != nil {
		log.Print("error sending initial state over websocket connection: ", err)
		return
	}

	m.connMu.Lock()
	m.conn = conn
	m.connMu.Unlock()
}

func (m *Master) getState() State {
	m.mu.Lock()

	m.lastUpdate = time.Now()
	jobID := m.Job.ID
	tasks := make(map[int]api.TaskState, len(m.Tasks))

	for id, task := range m.Tasks {
		tasks[id] = *task
	}

	m.mu.Unlock()

	// send encoded image and current generation
	return State{
		JobID:     jobID,
		StartedAt: m.Job.StartedAt,
		Tasks:     tasks,
	}
}

func (m *Master) sendOutput(output *gg.Context, generation uint, fitness float64) {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.conn == nil {
		return
	}

	state := m.getState()
	state.Generation = generation

	if img, err := util.EncodeImage(output.Image()); err == nil {
		state.Output = img
	}

	if fitness != 0 {
		state.Fitness = 1.0 / fitness
	}

	if err := m.conn.WriteJSON(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}

func (m *Master) sendUpdate() {
	state := m.getState()
	if err := m.conn.WriteJSON(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}

func (m *Master) sendData(state State) error {
	m.connMu.Lock()
	defer m.connMu.Unlock()

	if m.conn == nil {
		return nil
	}

	return m.conn.WriteJSON(state)
}

func (m *Master) sendEdges() {
	log.Printf("[task-generator] sending edges")

	m.mu.Lock()
	state := State{JobID: m.Job.ID, TargetImageEdges: m.TargetImageEdges}
	m.mu.Unlock()

	if err := m.sendData(state); err != nil {
		log.Printf("[task-generator] error sending edges: %v", err)
	}
}

func (m *Master) sendPalette() {
	log.Printf("[task-generator] sending palette")

	m.mu.Lock()
	state := State{JobID: m.Job.ID, Palette: m.Palette}
	m.mu.Unlock()

	if err := m.sendData(state); err != nil {
		log.Printf("[task-generator] error sending palette: %v", err)
	}
}
