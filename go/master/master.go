package master

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/db"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Master struct {
	Fitness           float64
	Generation        uint
	Job               api.Job
	NumWorkers        int
	Palette           string
	TargetImage       util.Image
	TargetImageBase64 string
	TargetImageEdges  string
	Tasks             map[int]*api.TaskState
	ThreadsPerWorker  int

	db                 db.DB
	conn               *websocket.Conn
	connMu             sync.Mutex
	lastUpdate         time.Time
	mu                 sync.Mutex
	transitioning      bool
	wsHeartbeatTimeout time.Duration
}

func Run() {
	numWorkers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		log.Fatal("error parsing WORKERS: ", err)
	}

	workerThreads, err := strconv.Atoi(os.Getenv("WORKER_THREADS"))
	if err != nil {
		log.Fatal("error parsing WORKER_THREADS: ", err)
	}

	m := Master{
		db:         db.NewConnection(),
		NumWorkers: numWorkers,
		Job: api.Job{
			ID:           1,
			CrossRate:    0.2,
			DetectEdges:  true,
			MutationRate: 0.022,
			NumColors:    64,
			NumShapes:    7000,
			OverDraw:     20,
			PaletteType:  "targetImage",
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    20,
			ShapeType:    "polygons",
		},
		lastUpdate:         time.Now(),
		TargetImage:        util.Image{},
		Tasks:              map[int]*api.TaskState{},
		ThreadsPerWorker:   workerThreads,
		wsHeartbeatTimeout: 2 * time.Second,
	}

	if !m.restoreFromSnapshot() && os.Getenv("START_RANDOM_JOB") == "true" {
		m.startRandomJob()
	}

	go m.httpServer()

	go m.detectFailures()

	go m.combine()

	go m.saveSnapshots()

	m.rpcServer()
}
