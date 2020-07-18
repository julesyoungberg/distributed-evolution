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
	NumWorkers       int                    `json:"numWorkers"`
	Output           string                 `json:"output"`
	TargetImage      string                 `json:"targetImage"`
	Tasks            map[int]*api.TaskState `json:"tasks"`
	ThreadsPerWorker int                    `json:"threadsPerWorker"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// keep websocket connection alive for as long as possible
// by periodically sending a message
func (m *Master) keepAlive(c *websocket.Conn) {
	for {
		m.connMu.Lock()

		if time.Since(m.lastUpdate) > m.wsHeartbeatTimeout {
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				m.connMu.Unlock()
				return
			}

			m.lastUpdate = time.Now()
		}

		m.connMu.Unlock()

		time.Sleep(200 * time.Millisecond)
	}
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

	m.connMu.Lock()

	log.Printf("new websocket connection request")

	// save the connection for updates
	m.conn = conn

	response := State{TargetImage: m.TargetImageBase64}
	if err := conn.WriteJSON(response); err != nil {
		log.Print("error sending initial state over websocket connection: ", err)
	}

	m.connMu.Unlock()

	go m.keepAlive(conn)
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
		Output:     img,
		Tasks:      tasks,
	}

	if err := m.conn.WriteJSON(state); err != nil {
		log.Print("[combiner] error sending output: ", err)
	}
}
