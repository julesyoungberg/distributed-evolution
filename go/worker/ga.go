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
func (w *Worker) createCallback(id uint32, thread int) func(ga *eaopt.GA) {
	// send the currrent best fit and other data to the master
	return func(ga *eaopt.GA) {
		w.mu.Lock()
		state := w.Tasks[id]
		w.mu.Unlock()

		// get best fit
		bestFit := ga.HallOfFame[0]

		output := state.BestFit.Output

		if output == nil {
			log.Printf("error! best fit image is nil at generation %v - bestFit: %v", ga.Generations, state.BestFit)
			return
		}

		state.Task.Output = util.EncodeImage(output)
		bestFit.Genome = api.Shapes{}

		// clear state
		state.BestFit = api.Output{}

		// add data to the task
		state.Task.BestFit = bestFit
		state.Task.Generation = ga.Generations

		w.saveTaskSnapshot(state)

		task, err := state.Task.UpdateMaster("inprogress")
		if err != nil {
			log.Print("error ", err)
		}

		// if the master responded with a different job id we are out of date
		if state.Task.Job.ID != task.Job.ID {
			log.Printf("out of date job of %v, updating to %v", state.Task.Job.ID, task.Job.ID)
			state.Task.Job.ID = task.Job.ID
		}
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(taskID uint32, jobID uint32) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		return w.Tasks[taskID].Task.Job.ID != jobID
	}
}
