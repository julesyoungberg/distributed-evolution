package worker

import (
	"fmt"
	"image/color"
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func getTaskContext(task api.Task, palette []color.RGBA) (*api.TaskContext, error) {
	t := api.TaskContext{
		GenOffset: task.Generation,
		Palette:   palette,
		Task:      task,
	}

	// decode and save target image
	img, err := util.DecodeImage(task.TargetImage)
	if err != nil {
		return nil, fmt.Errorf("error decoding task target image: %v", err)
	}

	width, height := util.GetImageDimensions(img)
	t.TargetImage = util.Image{
		Image:  img,
		Width:  width,
		Height: height,
	}

	if len(task.Edges) > 0 {
		edges, err := util.DecodeImage(task.Edges)
		if err != nil {
			return nil, fmt.Errorf("error decoding task target image edges: %v", err)
		}

		t.Edges = edges
	}

	return &t, nil
}

// RunTask executes the genetic algorithm for a given task
func (w *Worker) RunTask(task api.Task, thread int) {
	log.Printf("[thread %v] assigned task %v of job %v with population len: %v", thread, task.ID, task.Job.ID, len(task.Population))

	t, err := getTaskContext(task, w.Palette)

	w.ga = CreateGA(task.Job)

	// create closure functions with context
	w.ga.Callback = w.createCallback(task.ID, thread)
	w.ga.EarlyStop = w.createEarlyStop(task.ID)
	factory := api.GetShapesFactory(t, task.Population)

	// t.Task.Population = eaopt.Individuals{}

	w.mu.Lock()
	w.Tasks[task.ID] = t
	w.mu.Unlock()

	// evolve
	err = w.ga.Minimize(factory)
	if err != nil {
		log.Print(err)
	}
}

func (w *Worker) updateMaster(state *api.TaskContext, thread int, status string) bool {
	err := state.Task.UpdateMaster(w.ID, thread, status)
	if err != nil {
		log.Printf("[thread %v] failed to update master: %v", thread, err)
		state.Task.ID = -1
		return false
	}

	return true
}

// saves a task snapshot as a serialized JSON string to the cache
func (w *Worker) saveTaskSnapshot(task api.Task, thread int) {
	err := w.db.SaveTask(task)
	if err != nil {
		log.Printf("[thread %v] error saving task %v snapshot: %v", thread, task.ID, err)
	}
}

func (w *Worker) updateTask(state *api.TaskContext, ga *eaopt.GA, thread int) {
	// enrich task with data to save in the db
	task, err := state.EnrichTask(ga)
	if err != nil {
		log.Printf("[thread %v] error updateing task: %v", thread, err)
		return
	}

	// clear best fit
	state.BestFit = api.Output{}

	go w.saveTaskSnapshot(task, thread)
}
