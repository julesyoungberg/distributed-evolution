package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
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

		state.Task.Generation = state.GenOffset + ga.Generations

		success := w.updateMaster(state, thread)
		if !success {
			return
		}

		go func() {
			w.updateTask(state, ga, thread)

			w.mu.Lock()
			w.Tasks[id] = state
			w.mu.Unlock()
		}()
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(taskID int, jobID int) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		w.mu.Lock()
		jID := w.Tasks[taskID].Task.Job.ID
		w.mu.Unlock()

		return jID != jobID
	}
}
