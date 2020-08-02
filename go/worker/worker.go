package worker

import (
	"image/color"
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
	Palette      []color.RGBA
	Tasks        map[int]*api.WorkerTask

	db db.DB
	ga *eaopt.GA
	mu sync.Mutex
}

func (w *Worker) getPalette() error {
	palette, err := w.db.GetPalette()
	if err != nil {
		return err
	}

	w.mu.Lock()
	w.Palette = palette
	w.mu.Unlock()

	return nil
}

func (w *Worker) thread(thread int) {
	for {
		time.Sleep(10 * time.Second)

		if len(w.Palette) == 0 {
			err := w.getPalette()
			if err != nil {
				log.Printf("[thread %v] error getting palette: %v", thread, err)
				continue
			}
		}

		log.Printf("[thread %v] getting task", thread)

		task, err := w.db.PullTask()
		if err != nil {
			log.Printf("[thread %v] error: %v", thread, err)
			continue
		}

		if task.Generation != 0 {
			w.RunTask(task, thread)
			log.Printf("[thread %v] finished task %v", thread, task.ID)
			w.Palette = []color.RGBA{} // clear the palette
		}
	}
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
		db:           db.NewConnection(),
		NGenerations: 1000000000000, // 1 trillion
		Tasks:        map[int]*api.WorkerTask{},
	}

	nThreads, err := strconv.Atoi(os.Getenv("THREADS"))
	if err != nil {
		log.Fatalf("invalid THREADS value: %v", err)
	}

	// wait for a palette to be set
	for len(w.Palette) == 0 {
		if err = w.getPalette(); err != nil {
			time.Sleep(time.Second)
		}
	}

	log.Print("threads: ", nThreads)

	var wg sync.WaitGroup

	for i := 1; i <= nThreads; i++ {
		wg.Add(1)
		go w.thread(i + 1)
		time.Sleep(time.Second) // stagger threads to stagger requests to db
	}

	wg.Wait()
}
