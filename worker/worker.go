package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"strings"
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

func (w *Worker) runTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)

	w.currentTask = task

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(task.TargetImage))
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal("error decoding task target image ", err)
	}

	bounds := img.Bounds()
	w.targetImage = Image{
		image:  img,
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}

	Factory := createTriangleFactory(w)

	err = w.ga.Minimize(Factory)
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
		task := getTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.runTask(task)
		}

		time.Sleep(time.Second)
	}
}
