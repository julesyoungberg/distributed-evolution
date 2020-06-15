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
	gaConfig := eaopt.NewDefaultGAConfig()

	gaConfig.NPops = 1
	gaConfig.NGenerations = 1

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}

func (w *Worker) runTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)

	w.currentTask = task

	img := util.DecodeImage(task.TargetImage)

	bounds := img.Bounds()
	w.targetImage = Image{
		image:  img,
		width:  bounds.Dx(),
		height: bounds.Dy(),
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

	w.ga.Callback = func(ga *eaopt.GA) {
		fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)
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
