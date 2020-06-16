package main

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
	image  image.Image
	width  int
	height int
}

type Worker struct {
	ga           *eaopt.GA
	mutationRate float64
	currentTask  api.Task
	targetImage  Image
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
func (w *Worker) runTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)

	img := util.DecodeImage(task.TargetImage)
	width, height := util.GetImageDimensions(img)

	w.currentTask = task
	w.targetImage = Image{
		image:  img,
		width:  width,
		height: height,
	}

	w.ga.NGenerations = 1000

	w.ga.Callback = func(ga *eaopt.GA) {
		// TODO send current population to the master
		fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)
	}

	Factory := createTriangleFactory(w)

	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	w := Worker{
		ga:           createGA(),
		mutationRate: 0.8,
	}

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for {
		task := api.GetTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.runTask(task)
		}

		time.Sleep(time.Second)
	}
}
