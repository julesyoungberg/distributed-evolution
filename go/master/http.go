package master

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func (m *Master) respondWithState(w http.ResponseWriter) {
	m.mu.Lock()

	response := State{
		NumWorkers:       m.NumWorkers,
		Tasks:            m.Tasks,
		ThreadsPerWorker: m.ThreadsPerWorker,
	}

	m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handler for POST /job requests
// abandons current job and start on the new one
func (m *Master) newJob(w http.ResponseWriter, r *http.Request) {
	cors(w)

	// ignore preflight request
	if r.Method == http.MethodOptions {
		return
	}

	log.Print("##### new job request #####")

	// decode request body as job config
	var job api.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		log.Printf("error decoding new job request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decode and save target image
	base64 := job.TargetImage
	img, err := util.DecodeImage(job.TargetImage)
	if err != nil {
		log.Printf("error decoding target image")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.respondWithState(w)

	m.mu.Lock()

	m.TargetImageBase64 = base64
	m.setTargetImage(img)

	// save the job with a new ID
	newID := m.Job.ID + 1
	m.Job = job
	m.Job.ID = newID
	m.Job.TargetImage = "" // no need to be passing it around, its saved on m

	// clear the queue of all tasks
	for {
		_, err := m.db.PullTask()
		if err != nil {
			break
		}
	}

	// mark any qeued tasks as stale
	for _, task := range m.Tasks {
		if task.Status == "queued" {
			task.Status = "stale"
		}
	}

	// wait for all the workers to stop
	for !m.allStale() {
		m.mu.Unlock()
		time.Sleep(time.Second)
		m.mu.Lock()
	}

	m.mu.Unlock()

	m.generateTasks()
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, `{"alive": true}`); err != nil {
		util.DPrintf("health check error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Master) getKeyFromRedis(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	data, err := m.db.Get(params["key"])
	if err != nil {
		e := fmt.Sprintf("error getting key %v from redis: %v", params["key"], err)
		util.DPrintf(e)
		w.WriteHeader(http.StatusInternalServerError)

		if _, err := io.WriteString(w, e); err != nil {
			util.DPrintf("error writing error to response: %v", err)
		}
	}

	if _, err := io.WriteString(w, data); err != nil {
		util.DPrintf("error writing value from redis key %v to response: %v", params["key"], err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

// handles requests from the ui and websocket communication
func (m *Master) httpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/job", m.newJob).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/subscribe", m.subscribe)
	r.HandleFunc("/api/healthz", healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/api/redis/{key:[0-9A-Za-z:]+}", m.getKeyFromRedis).Methods(http.MethodGet)

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	r.Use(mux.CORSMethodMiddleware(r))

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	log.Fatal(http.ListenAndServe(":"+port, loggedRouter))
}
