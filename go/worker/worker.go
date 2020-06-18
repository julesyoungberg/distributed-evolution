package worker

import (
	"fmt"
	"image"
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

// TODO figure out how we can set an initial population to start from
// maybe make another version of createTriangleFactory that accepts a seed population
func (w *Worker) RunTask(task api.Task) {
	util.DPrintf("assigned task %v\n", task.ID)

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

	w.ga.NGenerations = 100
	w.ga.PopSize = 10

	w.ga.Callback = createCallback(task)

	Factory := createTriangleFactory(w)

	util.DPrintf("evolving...")

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
		// TODO handle errors by waiting and trying again
		task := api.GetTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.RunTask(task)
		}

		time.Sleep(time.Second)
	}
}
