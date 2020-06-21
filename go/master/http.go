package master

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/MaxHalford/eaopt"
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

	if r.Method == http.MethodOptions {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("starting new job")

	var job api.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("decoding target image")

	m.TargetImage = util.DecodeImage(job.TargetImage)
	prev := m.TargetImageBase64
	m.TargetImageBase64 = job.TargetImage

	log.Printf("SANITY CHECK 1 - job.TargetImage != prev: %v", job.TargetImage == prev)
	log.Printf("SANITY CHECK 2 - prevTarget != newTarget: %v", prev == m.TargetImageBase64)

	if job.TargetImage == prev {
		log.Printf("error, new target image is the same as current")
	}

	m.Job = job
	m.Job.ID = uuid.New().ID()
	m.Job.TargetImage = "" // no need to be passing it around, its saved on m

	m.Generations = Generations{}
	m.Outputs = map[int]Generation{}

	log.Printf("generating tasks")

	m.generateTasks()

	state := State{
		NumWorkers: m.NumWorkers,
		Tasks:      make([]api.Task, len(m.Tasks)),
	}

	for i, t := range m.Tasks {
		t.BestFit = eaopt.Individual{}
		state.Tasks[i] = t
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// keep websocket connection alive for as long as possible
func (m *Master) keepAlive(c *websocket.Conn) {
	go func() {
		for {
			m.mu.Lock()

			if time.Since(m.lastUpdate) > time.Second {
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
	}()
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

	m.conn = conn

	state := State{TargetImage: m.TargetImageBase64}

	log.Printf("sending data")

	if err := conn.WriteJSON(state); err != nil {
		log.Println(err)
	}

	m.mu.Unlock()

	m.keepAlive(conn)
}

func (m *Master) sendOutput(output *gg.Context) {
	if m.conn == nil {
		util.DPrintf("no open ui connections")
		return
	}

	// get resulting image
	img := output.Image()

	// send encoded image and current generation
	state := State{
		Output: util.EncodeImage(img),
		Tasks:  make([]api.Task, len(m.Tasks)),
	}

	var latest uint = 0

	for i, t := range m.Tasks {
		t.BestFit = eaopt.Individual{}
		state.Tasks[i] = t

		if t.Generation > latest {
			latest = t.Generation
		}
	}

	state.Generation = latest

	util.DPrintf("sending generation %v update to ui", latest)

	if err := m.conn.WriteJSON(state); err != nil {
		log.Println(err)
	}

	m.lastUpdate = time.Now()
}

// sends the given generation's output image to the UI
func (m *Master) updateUICombined(generation Generation) {
	genN := generation.Generation

	util.DPrintf("updating ui with generation %v", genN)
	util.DPrintf("encoding output image for generation %v", genN)

	m.sendOutput(generation.Output)
}

// draws the latest generations to a single image
func (m *Master) updateUILatest() {
	dc := gg.NewContext(m.TargetImageWidth, m.TargetImageHeight)

	for _, t := range m.Tasks {
		// if t.BestFit.Genome == nil {
		// 	continue
		// }

		// s := t.BestFit.Genome.(worker.Shapes)
		// s.Draw(dc, t.Offset)

		out, ok := m.Outputs[t.ID]
		if !ok {
			continue
		}

		centerX := int(math.Round(t.Offset.X + t.Dimensions.X/2.0))
		centerY := int(math.Round(t.Offset.Y + t.Dimensions.Y/2.0))

		dc.DrawImageAnchored(out.Image, centerX, centerY, 0.5, 0.5)
	}

	m.sendOutput(dc)
}

// TODO take target image from http
// allow multiple jobs at the same time, take number of workers from request too??
func (m *Master) httpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/job", m.newJob).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/subscribe", m.subscribe)
	// http.Handle("/subscribe", websocket.Handler(m.subscribe))

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe(":"+port, r))
}
