package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
)

// initialize GA
func CreateGA(config api.Job) *eaopt.GA {
	nGenerations := config.NumGenerations
	if nGenerations < 1 {
		nGenerations = 1000000000000 // 1 trillion
	}

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

func (w *Worker) setTaskState(state *api.WorkerTask) {
	w.mu.Lock()
	w.Tasks[state.Task.ID] = state
	w.mu.Unlock()
}

// returns a closure to be called after each generation
func (w *Worker) createCallback(id int, thread int) func(ga *eaopt.GA) {
	// send the currrent best fit and other data to the master
	return func(ga *eaopt.GA) {
		w.mu.Lock()
		state := w.Tasks[id]
		w.mu.Unlock()

		generation := state.GenOffset + ga.Generations
		nGenerations := state.Task.Job.NumGenerations
		state.Task.Generation = generation

		status := "inprogress"
		if nGenerations > 0 && generation >= nGenerations {
			status = "done"
		}

		success := w.updateMaster(state, thread, status)
		if !success {
			w.setTaskState(state)
			return
		}

		go func() {
			w.updateTask(state, ga, thread)
			w.setTaskState(state)
		}()
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(taskID int) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		w.mu.Lock()
		state := w.Tasks[taskID]
		taskID := state.Task.ID
		generation := state.Task.Generation
		nGenerations := state.Task.Job.NumGenerations
		w.mu.Unlock()

		// extra check because eaopt seems to disregard
		done := nGenerations > 0 && generation >= nGenerations

		return taskID == -1 || done
	}
}
