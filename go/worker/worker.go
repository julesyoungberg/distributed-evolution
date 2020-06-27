package worker

import (
	"image"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/google/uuid"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/db"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Output struct {
	Fitness float64
	Output  image.Image
}

type WorkerTask struct {
	BestFit     Output
	TargetImage util.Image
	Task        api.Task

	mu sync.Mutex
}

type Worker struct {
	ID           uint32
	NGenerations uint
	Tasks        map[uint32]*WorkerTask

	db db.DB
	ga *eaopt.GA
	mu sync.Mutex
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
		db:           db.NewConnection(),
		NGenerations: 1000000000000, // 1 trillion
		Tasks:        map[uint32]*WorkerTask{},
	}

	nThreads, err := strconv.Atoi(os.Getenv("THREADS"))
	if err != nil {
		log.Fatalf("invalid THREADS value: %v", err)
	}

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for i := 0; i < nThreads; i++ {
		go func(thread int) {
			for {
				task, err := w.db.PullTask()

				if err == nil && task.Generation != 0 {
					log.Print("assigned task ", task.ID)
					w.RunTask(task, thread)
					log.Print("finished task ", task.ID)
				} else if err != nil {
					// TODO change to log.Print for fault tolerance
					log.Print("error: ", err)
				} else {
					log.Print("empty task response, waiting for work")
				}

				time.Sleep(10 * time.Second)
			}
		}(i)
	}
}
