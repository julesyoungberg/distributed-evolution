package master

import (
	"encoding/json"
	"image"
	"image/color"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/cv"
	"github.com/rickyfitts/distributed-evolution/go/util"
	"github.com/rickyfitts/distributed-evolution/go/worker"
)

type SingleSystem struct {
	Master           *Master
	Output           image.Image
	Palette          []color.RGBA
	TargetImageEdges image.Image
	Task             api.Task
	WorkerTask       *api.WorkerTask

	ga *eaopt.GA
}

func newSingleSystem() SingleSystem {
	m := newMaster()
	return SingleSystem{Master: &m}
}

func (s *SingleSystem) preparePalette() {
	log.Print("[task-generator] preparing palette")

	palette := s.Master.getPalette()

	s.Master.mu.Lock()
	s.Palette = palette
	s.Master.mu.Unlock()

	s.Master.savePalete(s.Palette)
}

func (s *SingleSystem) generateTask() {
	log.Print("[task-generator] generating task")

	s.Master.mu.Lock()

	s.Master.Job.ShapesPerSlice = s.Master.Job.NumShapes
	job := s.Master.Job

	targetImage := s.Master.TargetImage.Image
	width := float64(s.Master.TargetImage.Width)
	height := float64(s.Master.TargetImage.Height)

	s.Master.mu.Unlock()

	s.preparePalette()

	if job.DetectEdges {
		log.Printf("[task-generator] getting target image edges")
		edges, err := cv.GetEdges(targetImage)
		if err != nil {
			log.Fatal(err)
		}

		go s.Master.saveEdges(edges)

		s.Master.mu.Lock()
		s.TargetImageEdges = edges
		s.Master.mu.Unlock()
	} else {
		s.Master.mu.Lock()
		s.TargetImageEdges = nil
		s.Master.mu.Unlock()
	}

	task := api.Task{
		Dimensions:         util.Vector{X: width, Y: height},
		ID:                 1,
		Job:                job,
		ScaledQuantization: job.Quantization,
		ShapeType:          job.ShapeType,
	}

	s.Master.mu.Lock()
	s.Task = task
	s.Master.mu.Unlock()
}

func (s *SingleSystem) createCallback() func(ga *eaopt.GA) {
	return func(ga *eaopt.GA) {
		s.Master.mu.Lock()
		defer s.Master.mu.Unlock()

		state := s.WorkerTask
		output := state.BestFit.Output
		fitness := state.BestFit.Fitness

		if fitness != 0 {
			fitness = 1 / fitness
		}

		s.Master.Fitness = fitness
		s.Master.Generation = ga.Generations
		s.Output = output

		generation := s.Master.Generation
		nGenerations := s.Master.Job.NumGenerations

		if nGenerations > 0 && generation >= nGenerations {
			s.Master.Job.Complete = true
			s.Master.Job.CompletedAt = time.Now()
		}
	}
}

func (s *SingleSystem) createEarlyStop() func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		s.Master.mu.Lock()
		state := s.WorkerTask
		s.Master.mu.Unlock()

		state.Mu.Lock()
		nGenerations := state.Task.Job.NumGenerations
		generation := state.Task.Generation
		state.Mu.Unlock()

		// extra check because eaopt seems to disregard
		return nGenerations > 0 && generation >= nGenerations
	}
}

func (s *SingleSystem) runTask() {
	log.Printf("[task-runner] running task - quantization: %v", s.Task.Job.Quantization)

	s.Master.mu.Lock()
	job := s.Master.Job
	task := s.Task
	t := api.WorkerTask{
		GenOffset:   0,
		Palette:     s.Palette,
		TargetImage: s.Master.TargetImage,
		Task:        task,
	}
	s.Master.mu.Unlock()

	if job.DetectEdges {
		t.Edges = s.TargetImageEdges
	}

	s.ga = worker.CreateGA(job)
	s.ga.Callback = s.createCallback()
	s.ga.EarlyStop = s.createEarlyStop()
	factory := api.GetShapesFactory(&t, task.Population)

	s.Master.mu.Lock()
	s.WorkerTask = &t
	s.Master.mu.Unlock()

	// evolve
	err := s.ga.Minimize(factory)
	if err != nil {
		log.Print(err)
	}
}

func (s *SingleSystem) startRandomJob() {
	s.Master.setRandomTargetImage()
	s.Master.Job.StartedAt = time.Now()
	s.generateTask()
	go s.runTask()
}

func (s *SingleSystem) updateUI() {
	for {
		time.Sleep(5 * time.Second)

		s.Master.mu.Lock()

		output := s.Output
		generation := s.Master.Generation
		fitness := s.Master.Fitness

		s.Master.mu.Unlock()

		if output == nil || generation < 1 {
			log.Printf("output nil or generation < 1")
			continue
		}

		log.Printf("updating UI (generation: %v, fitness: %v)", generation, fitness)

		s.Master.sendOutput(output)
	}
}

func (s *SingleSystem) newJob(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("error decoding new job request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decode and save target image
	img, err := util.DecodeImage(job.TargetImage)
	if err != nil {
		log.Printf("error decoding target image: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.Master.mu.Lock()

	// save the job with a new ID
	job.ID = s.Master.Job.ID + 1
	s.Master.Job = job

	s.Master.TargetImageBase64 = job.TargetImage
	s.Master.setTargetImage(img)
	s.Master.Job.TargetImage = "" // no need to be passing it around, its saved on m
	s.Master.Job.StartedAt = time.Now()
	s.Master.Job.Complete = false
	s.Master.Generation = 0
	s.Master.Fitness = 0.0

	s.Master.mu.Unlock()

	s.generateTask()
	go s.runTask()

	s.Master.respondWithState(w)
}

func (s *SingleSystem) httpServer() {
	r := mux.NewRouter()

	r.HandleFunc("/api/job", s.newJob).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/subscribe", s.Master.subscribe)
	r.HandleFunc("/api/healthz", healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/api/palette", s.Master.fetchPalette).Methods(http.MethodGet)
	r.HandleFunc("/api/state", s.Master.fetchState).Methods(http.MethodGet)

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	r.Use(mux.CORSMethodMiddleware(r))

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	log.Fatal(http.ListenAndServe(":"+port, loggedRouter))
}

func RunSingleSystem() {
	log.Print("starting single system")

	s := newSingleSystem()

	if os.Getenv("START_RANDOM_JOB") == "true" {
		s.startRandomJob()
	}

	go s.updateUI()

	s.httpServer()
}
