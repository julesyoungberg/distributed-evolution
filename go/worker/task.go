package task

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

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
