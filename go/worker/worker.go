package worker

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/google/uuid"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/db"
)

type Worker struct {
	ID           uint32
	NGenerations uint
	Tasks        map[uint32]*api.WorkerTask

	db db.DB
	ga *eaopt.GA
	mu sync.Mutex
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
		db:           db.NewConnection(),
		NGenerations: 1000000000000, // 1 trillion
		Tasks:        map[uint32]*api.WorkerTask{},
	}

	nThreads, err := strconv.Atoi(os.Getenv("THREADS"))
	if err != nil {
		log.Fatalf("invalid THREADS value: %v", err)
	}

	log.Print("threads: ", nThreads)

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	var wg sync.WaitGroup

	for i := 0; i < nThreads; i++ {
		wg.Add(1)

		go func(thread int) {
			for {
				log.Printf("[thread %v] getting task", thread)

				task, err := w.db.PullTask()

				if err == nil && task.Generation != 0 {
					w.RunTask(task, thread)
					log.Printf("[thread %v] finished task %v", thread, task.ID)
				} else if err != nil {
					log.Printf("[thread %v] error: %v", thread, err)
				} else {
					log.Printf("[thread %v] empty task response, waiting for work", thread)
				}

				time.Sleep(10 * time.Second)
			}
		}(i)
	}

	wg.Wait()
}
