package worker

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/MaxHalford/eaopt"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Image struct {
	Image  image.Image
	Width  int
	Height int
}

type Worker struct {
	MutationRate float64
	CurrentTask  api.Task
	TargetImage  Image

	ga *eaopt.GA
}

func createGA() *eaopt.GA {
	ga, err := eaopt.NewDefaultGAConfig().NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}

// TODO figure out how we can set an initial population to start from
// maybe make another version of createTriangleFactory that accepts a seed population
func (w *Worker) RunTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)

	img := util.DecodeImage(task.TargetImage)
	width, height := util.GetImageDimensions(img)

	w.CurrentTask = task
	w.TargetImage = Image{
		Image:  img,
		Width:  width,
		Height: height,
	}

	w.ga.NGenerations = 1000

	w.ga.Callback = func(ga *eaopt.GA) {
		fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)

		task.Population = ga.Populations[0]

		api.Update(task)
	}

	Factory := createTriangleFactory(w)

	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}
}

func Run() {
	w := Worker{MutationRate: 0.8}
	w.ga = createGA()

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for {
		task := api.GetTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.RunTask(task)
		}

		time.Sleep(time.Second)
	}
}
