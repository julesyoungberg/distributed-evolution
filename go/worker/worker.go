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

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) saveTaskSnapshot(state *WorkerTask) {
	task := state.Task
	task.Population = w.ga.Populations[0]
	w.db.SaveTask(task)
}

// RunTask executes the genetic algorithm for a given task
func (w *Worker) RunTask(task api.Task, thread int) {
	population := task.Population

	task.Population = eaopt.Population{}
	task.WorkerID = w.ID
	task.Thread = thread

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

	// create closure functions with context
	w.ga.Callback = w.createCallback(task.ID, thread)
	w.ga.EarlyStop = w.createEarlyStop(task.ID, task.Job.ID)
	factory := getShapesFactory(&t, population)

	w.Tasks[task.ID] = &t

	task.UpdateMaster("inprogress")

	// evolve
	err := w.ga.Minimize(factory)
	if err != nil {
		log.Print(err)
	}
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
