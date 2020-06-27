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
	Generation       uint              `json:"generation"`
	NumWorkers       int               `json:"numWorkers"`
	Output           string            `json:"output"`
	TargetImage      string            `json:"targetImage"`
	Tasks            map[int]*api.Task `json:"tasks"`
	ThreadsPerWorker int               `json:"threadsPerWorker"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// keep websocket connection alive for as long as possible
// by periodically sending a message
func (m *Master) keepAlive(c *websocket.Conn) {
	for {
		m.mu.Lock()

		if time.Since(m.lastUpdate) > m.wsHeartbeatTimeout {
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				m.mu.Unlock()
				return
			}

			m.lastUpdate = time.Now()
		}

		m.mu.Unlock()

		time.Sleep(200 * time.Millisecond)
	}
}

// subscribe handler creates a websocket connection with the client
// TODO create connection mutex not to block other tasks
// TODO handle multiple connections
func (m *Master) subscribe(w http.ResponseWriter, r *http.Request) {
	// TODO check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	m.mu.Lock()

	log.Printf("new websocket connection request")

	// save the connection for updates
	m.conn = conn

	response := State{TargetImage: m.TargetImageBase64}
	if err := conn.WriteJSON(response); err != nil {
		log.Println(err)
	}

	m.mu.Unlock()

	go m.keepAlive(conn)
}

func (m *Master) sendOutput(output *gg.Context, generation uint) {
	m.lastUpdate = time.Now()

	if m.conn == nil {
		// no open connections
		return
	}

	// get resulting image
	img, err := util.EncodeImage(output.Image())
	if err != nil {
		log.Print("error sending output: ", err)
		return
	}

	// send encoded image and current generation
	state := State{
		Generation: generation,
		Output:     img,
		Tasks:      m.Tasks,
	}

	if err := m.conn.WriteJSON(state); err != nil {
		log.Println(err)
	}

	m.lastUpdate = time.Now()
}