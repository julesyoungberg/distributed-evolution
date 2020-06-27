package master

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
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
	Tasks             map[uint32]*api.Task
	ThreadsPerWorker  int

	db                 db.DB
	conn               *websocket.Conn
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
		log.Fatal("error parsing WORKERS: ", err)
	}

	m := Master{
		db:         db.NewConnection(),
		NumWorkers: numWorkers,
		Job: api.Job{
			ID:           uuid.New().ID(),
			CrossRate:    0.2,
			MutationRate: 0.021,
			NumShapes:    200,
			OverDraw:     10,
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    20,
		},
		lastUpdate:         time.Now(),
		TargetImage:        util.Image{},
		Tasks:              map[uint32]*api.Task{},
		ThreadsPerWorker:   workerThreads,
		wsHeartbeatTimeout: 2 * time.Second,
	}

	log.Print("fetching random image...")
	image := util.GetRandomImage()

	log.Print("encoding image...")
	m.TargetImageBase64 = util.EncodeImage(image)

	m.setTargetImage(image)

	go m.generateTasks()

	go m.httpServer()

	go m.detectFailures()

	go m.combine()

	m.rpcServer()
}
