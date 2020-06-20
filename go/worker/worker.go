package worker

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/google/uuid"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Image struct {
	Image  image.Image
	Width  int
	Height int
}

type Worker struct {
	ID           uint32
	CurrentTask  api.Task
	Job          api.Job
	NGenerations uint
	TargetImage  Image

	ga *eaopt.GA
}

// TODO figure out how we can set an initial population to start from
// maybe make another version of createTriangleFactory that accepts a seed population
func (w *Worker) RunTask(task api.Task) {
	log.Printf("assigned task %v\n", task.ID)

	util.DPrintf("decoding image...")
	img := util.DecodeImage(task.TargetImage)
	width, height := util.GetImageDimensions(img)

	util.DPrintf("saving task data...")
	w.CurrentTask = task
	w.TargetImage = Image{
		Image:  img,
		Width:  width,
		Height: height,
	}

	util.DPrintf("preparing ga...")
	util.DPrintf("setting ID to %v", task.Job.ID)

	w.Job = task.Job
	util.DPrintf("ID: %v", w.Job.ID)

	w.createGA()

	task.Job = api.Job{ID: task.Job.ID}

	w.ga.Callback = w.createCallback(task)
	w.ga.EarlyStop = w.createEarlyStop(task)
	Factory := createShapesFactory(w, task.Type)

	util.DPrintf("evolving...")

	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("finishing task")
}

func Run() {
	w := Worker{
		ID:           uuid.New().ID(),
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
		// if generation is zero this is an empty response, if so just wait for more work
		if err == nil && task.Generation != 0 {
			w.RunTask(task)
		} else if err != nil {
			log.Printf("error: %v", err)
		} else {
			log.Print("empty task response...")
		}

		time.Sleep(time.Second)
	}
}
