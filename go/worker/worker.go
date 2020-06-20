package worker

import (
	"fmt"
	"image"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Image struct {
	Image  image.Image
	Width  int
	Height int
}

type Worker struct {
	CrossRate    float64
	CurrentTask  api.Task
	MutationRate float64
	NGenerations uint
	NumShapes    int
	PoolSize     uint
	PopSize      uint
	ShapeSize    float64
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

	w.ga.Callback = createCallback(task)

	Factory := createShapesFactory(w, task.Type)

	util.DPrintf("evolving...")

	err := w.ga.Minimize(Factory)
	if err != nil {
		fmt.Println(err)
	}
}

func Run() {
	w := Worker{
		CrossRate:    0.2,
		MutationRate: 0.021,
		NGenerations: 100000,
		NumShapes:    100,
		PoolSize:     20,
		PopSize:      100,
		ShapeSize:    50,
	}

	w.createGA()

	util.DPrintf("waiting for master to start")

	// wait for master to initialize
	time.Sleep(10 * time.Second)

	for {
		// TODO handle errors by waiting and trying again
		task := api.GetTask()

		// if generation is zero this is an empty response, if so just wait for more work
		if task.Generation != 0 {
			w.RunTask(task)
			break
		}

		time.Sleep(time.Second)
	}
}
