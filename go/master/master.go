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
	Tasks             map[int]*api.Task
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
		log.Fatal("error parsing WORKERS: ", err)
	}

	m := Master{
		db:         db.NewConnection(),
		NumWorkers: numWorkers,
		Job: api.Job{
			ID:           1,
			CrossRate:    0.2,
			MutationRate: 0.021,
			NumShapes:    300,
			OverDraw:     10,
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    30,
		},
		lastUpdate:         time.Now(),
		TargetImage:        util.Image{},
		Tasks:              map[int]*api.Task{},
		ThreadsPerWorker:   workerThreads,
		wsHeartbeatTimeout: 2 * time.Second,
	}

	log.Print("fetching random image...")
	image := util.GetRandomImage()

	log.Print("encoding image...")
	encodedImg, err := util.EncodeImage(image)
	if err != nil {
		log.Fatal(err)
	}

	m.TargetImageBase64 = encodedImg
	m.setTargetImage(image)

	go m.generateTasks()

	go m.httpServer()

	go m.detectFailures()

	go m.combine()

	m.rpcServer()
}
