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
	Job               api.Job
	NumWorkers        int
	TargetImage       util.Image
	TargetImageBase64 string
	Tasks             map[int]*api.TaskState
	ThreadsPerWorker  int

	db                 db.DB
	conn               *websocket.Conn
	connMu             sync.Mutex
	lastUpdate         time.Time
	mu                 sync.Mutex
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
			MutationRate: 0.021,
			NumShapes:    300,
			OverDraw:     20,
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    30,
		},
		lastUpdate:         time.Now(),
		TargetImage:        util.Image{},
		Tasks:              map[int]*api.TaskState{},
		ThreadsPerWorker:   workerThreads,
		wsHeartbeatTimeout: 2 * time.Second,
	}

	if !m.restoreFromSnapshot() {
		m.startRandomTask()
	}

	go m.httpServer()

	go m.detectFailures()

	go m.combine()

	go m.saveSnapshots()

	m.rpcServer()
}
