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
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Output struct {
	Fitness float64
	Output  image.Image
}

type Worker struct {
	ID           uint32
	BestFit      Output
	Job          api.Job
	NGenerations uint
	TargetImage  util.Image

	ga *eaopt.GA
	mu sync.Mutex
}

// RunTask executes the genetic algorithm for a given task
// TODO set an initial population to start from
func (w *Worker) RunTask(task api.Task) {
	// decode and save target image
	img := util.DecodeImage(task.TargetImage)
	width, height := util.GetImageDimensions(img)
	w.TargetImage = util.Image{
		Image:  img,
		Width:  width,
		Height: height,
	}

	// save job information for createGA to use
	w.Job = task.Job

	// clear job data from task to keep update messages small
	// the master only needs to confirm that the ID is correct
	task.Job = api.Job{ID: task.Job.ID}

	w.createGA()

	// create clsoure functions with context
	w.ga.Callback = w.createCallback(task)
	w.ga.EarlyStop = w.createEarlyStop(task)
	Factory := createShapesFactory(w, task.Type)

	// evolve
	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
		BestFit:      Output{},
		NGenerations: 1000000000000, // 1 trillion
		Job: api.Job{
			ID:           uuid.New().ID(),
			CrossRate:    0.2,
			MutationRate: 0.021,
			NumShapes:    200,
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    20,
		},
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
			log.Fatal("error: ", err)
		} else {
			log.Print("empty task response, waiting for work")
		}

		time.Sleep(time.Second)
	}
}
