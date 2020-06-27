package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) saveTaskSnapshot(state *api.WorkerTask, thread int) {
	task := state.Task
	task.Population = w.ga.Populations[0].Individuals
	err := w.db.SaveTask(task)
	if err != nil {
		log.Printf("[thread %v] error saving task %v snapshot: %v", thread, task.ID, err)
	}
}

// RunTask executes the genetic algorithm for a given task
func (w *Worker) RunTask(task api.Task, thread int) {
	log.Printf("[thread %v] assigned task %v", thread, task.ID)

	population := task.Population

	task.Population = eaopt.Individuals{}
	task.WorkerID = w.ID
	task.Thread = thread

	t := api.WorkerTask{Task: task}

	// decode and save target image
	img, err := util.DecodeImage(task.TargetImage)
	if err != nil {
		log.Printf("[thrad %v] error decoding task target image: %v", thread, err)
		return
	}

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
	factory := api.GetShapesFactory(&t, population)

	w.mu.Lock()
	w.Tasks[task.ID] = &t
	w.mu.Unlock()

	// evolve
	err = w.ga.Minimize(factory)
	if err != nil {
		log.Print(err)
	}
}
