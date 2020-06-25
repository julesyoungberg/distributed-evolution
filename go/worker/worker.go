package worker

import (
	"fmt"
	"image"
	"log"
	"sync"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/google/uuid"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/cache"
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
	Tasks        map[int]*WorkerTask

	cache cache.Cache
	ga    *eaopt.GA
	mu    sync.Mutex
}

// RunTask executes the genetic algorithm for a given task
// TODO set an initial population to start from
func (w *Worker) RunTask(task api.Task) {
	t := WorkerTask{Task: task}
	// decode and save target image
	img := util.DecodeImage(task.TargetImage)
	width, height := util.GetImageDimensions(img)
	t.TargetImage = util.Image{
		Image:  img,
		Width:  width,
		Height: height,
	}

	w.ga = createGA(task.Job, w.NGenerations)

	// create clsoure functions with context
	w.ga.Callback = w.createCallback(task.ID)
	w.ga.EarlyStop = w.createEarlyStop(task.ID, task.Job.ID)
	Factory := createShapesFactory(&t, task.Type)

	w.Tasks[task.ID] = &t

	// evolve
	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
		cache:        cache.NewConnection(),
		NGenerations: 1000000000000, // 1 trillion
		Tasks:        map[int]*WorkerTask{},
	}

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for {
		task, err := api.GetTask(w.ID)

		if err == nil && task.Generation != 0 {
			log.Print("assigned task ", task.ID)
			w.RunTask(task)
			log.Print("finished task ", task.ID)
		} else if err != nil {
			// TODO change to log.Print for fault tolerance
			log.Fatal("error: ", err)
		} else {
			log.Print("empty task response, waiting for work")
		}

		time.Sleep(time.Second)
	}
}
