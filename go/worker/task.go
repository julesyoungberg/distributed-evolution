package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// RunTask executes the genetic algorithm for a given task
func (w *Worker) RunTask(task api.Task, thread int) {
	log.Printf("[thread %v] assigned task %v of job %v with population len: %v", thread, task.ID, task.Job.ID, len(task.Population))

	t := api.WorkerTask{
		GenOffset: task.Generation,
		Palette:   w.Palette,
		Task:      task,
	}

	// decode and save target image
	img, err := util.DecodeImage(task.TargetImage)
	if err != nil {
		log.Printf("[thread %v] error decoding task target image: %v", thread, err)
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
	w.ga.EarlyStop = w.createEarlyStop(task.ID)
	factory := api.GetShapesFactory(&t, task.Population)

	t.Task.Population = eaopt.Individuals{}

	w.mu.Lock()
	w.Tasks[task.ID] = &t
	w.mu.Unlock()

	// evolve
	err = w.ga.Minimize(factory)
	if err != nil {
		log.Print(err)
	}
}

func (w *Worker) updateMaster(state *api.WorkerTask, thread int) bool {
	err := state.Task.UpdateMaster(w.ID, thread, "inprogress")
	if err != nil {
		log.Printf("[thread %v] failed to update master: %v", thread, err)
		state.Task.ID = -1
		return false
	}

	return true
}

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) saveTaskSnapshot(state *api.WorkerTask, thread int) {
	task := state.Task
	task.Population = w.ga.Populations[0].Individuals
	err := w.db.SaveTask(task)
	if err != nil {
		log.Printf("[thread %v] error saving task %v snapshot: %v", thread, task.ID, err)
	}
}

func (w *Worker) updateTask(state *api.WorkerTask, ga *eaopt.GA, thread int) {
	// get best fit
	bestFit := ga.HallOfFame[0]

	output := state.BestFit.Output
	if output == nil {
		// this happens A LOT - idk why
		// log.Printf("[thread %v] error! best fit image is nil at generation %v - bestFit: %v", thread, ga.Generations, state.BestFit)
		return
	}

	encoded, err := util.EncodeImage(output)
	if err != nil {
		log.Printf("[thread %v] error saving task: %v", thread, err)
		return
	}

	state.Task.Output = encoded
	bestFit.Genome = api.Shapes{}

	// clear state
	state.BestFit = api.Output{}

	// add data to the task
	state.Task.BestFit = bestFit

	w.saveTaskSnapshot(state, thread)
}
