package master

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fogleman/gg"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type State struct {
	Generation  uint       `json:"generation"`
	NumWorkers  int        `json:"numWorkers"`
	Output      string     `json:"output"`
	TargetImage string     `json:"targetImage"`
	Tasks       []api.Task `json:"tasks"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// handler for POST /job requests
// abandons current job and start on the new one
func (m *Master) newJob(w http.ResponseWriter, r *http.Request) {
	log.Printf("##### New Job Request - %v #####", http.MethodOptions)
	// Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

	// ignore preflight request
	if r.Method == http.MethodOptions {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("starting new job")

	// decode request body as job config
	var job api.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decode and save target image
	m.TargetImage = util.Image{Image: util.DecodeImage(job.TargetImage)}
	m.TargetImageBase64 = job.TargetImage

	// save the job with a new ID
	m.Job = job
	m.Job.ID = uuid.New().ID()
	m.Job.TargetImage = "" // no need to be passing it around, its saved on m

	// reset state and generate tasks
	m.Generations = Generations{}
	m.generateTasks()

	response := State{
		NumWorkers: m.NumWorkers,
		Tasks:      append([]api.Task{}, m.Tasks...),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	img := output.Image()

	// send encoded image and current generation
	state := State{
		Generation: generation,
		Output:     util.EncodeImage(img),
		Tasks:      append([]api.Task{}, m.Tasks...),
	}

	if err := m.conn.WriteJSON(state); err != nil {
		log.Println(err)
	}

	m.lastUpdate = time.Now()
}

// sends the given generation's output image to the UI
func (m *Master) updateUICombined(generation Generation) {
	m.sendOutput(generation.Output, generation.Generation)
}

// draws the latest generations to a single image
func (m *Master) updateUILatest() {
	latest := m.getLatestGeneration()
	m.sendOutput(latest.Output, latest.Generation)
}

// handles requests from the ui and websocket communication
func (m *Master) httpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/job", m.newJob).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/subscribe", m.subscribe)

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe(":"+port, r))
}
