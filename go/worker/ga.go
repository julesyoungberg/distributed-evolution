package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// initialize GA
func createGA(config api.Job, nGenerations uint) *eaopt.GA {
	gaConfig := eaopt.GAConfig{
		NPops:        1,
		PopSize:      config.PopSize,
		HofSize:      1,
		NGenerations: nGenerations,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: config.PoolSize,
			},
			MutRate:   config.MutationRate,
			CrossRate: config.CrossRate,
		},
		ParallelEval: false,
	}

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}

// returns a closure to be called after each generation
func (w *Worker) createCallback(id int, thread int) func(ga *eaopt.GA) {
	// send the currrent best fit and other data to the master
	return func(ga *eaopt.GA) {
		w.mu.Lock()
		state := w.Tasks[id]
		w.mu.Unlock()

		// get best fit
		bestFit := ga.HallOfFame[0]

		output := state.BestFit.Output

		if output == nil {
			log.Printf("[thread %v] error! best fit image is nil at generation %v - bestFit: %v", thread, ga.Generations, state.BestFit)
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
		state.Task.Generation = ga.Generations

		w.saveTaskSnapshot(state, thread)

		task, err := state.Task.UpdateMaster("inprogress")
		if err != nil {
			log.Printf("[thread %v] failed to update master %v", thread, err)
			state.Task.Job.ID = 0
			return
		}

		// if the master responded with a different job id we are out of date
		if state.Task.Job.ID != task.Job.ID {
			log.Printf("[thrad %v] out of date job of %v, updating to %v", thread, state.Task.Job.ID, task.Job.ID)
			state.Task.Job.ID = task.Job.ID
		}

		w.mu.Lock()
		w.Tasks[id] = state
		w.mu.Unlock()
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(taskID int, jobID int) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		w.mu.Lock()
		id := w.Tasks[taskID].Task.Job.ID
		defer w.mu.Unlock()

		return id != jobID
	}
}
