package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// initialize GA
func (w *Worker) createGA() {
	gaConfig := eaopt.GAConfig{
		NPops:        1,
		PopSize:      w.Job.PopSize,
		HofSize:      1,
		NGenerations: w.NGenerations,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: w.Job.PoolSize,
			},
			MutRate:   w.Job.MutationRate,
			CrossRate: w.Job.CrossRate,
		},
		ParallelEval: false,
	}

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	w.ga = ga
}

// returns a closure to be called after each generation
func (w *Worker) createCallback(task api.Task) func(ga *eaopt.GA) {
	// send the currrent best fit and other data to the master
	return func(ga *eaopt.GA) {
		// get best fit
		bestFit := ga.HallOfFame[0]

		// if draw once is active add the output and clear the genome
		if w.Job.DrawOnce {
			w.mu.Lock()
			output := w.BestFit.Output
			w.mu.Unlock()

			if output == nil {
				log.Printf("error! best fit image is nil at generation %v - bestFit: %v", ga.Generations, .BestFit)
				return
			}

			task.Output = util.EncodeImage(output)
			bestFit.Genome = Shapes{}

			// clear state
			w.BestFit = Output{}
		} else {
			// otherwise just add the genome
			bestFit.Genome = bestFit.Genome.(Shapes).CloneForSending()
		}

		// add data to the task
		task.BestFit = bestFit
		task.Generation = ga.Generations

		// send results to master
		jobId, err := api.Update(task)
		if err != nil {
			log.Print("error ", err)
		}

		// if the master responded with a different job id we are out of date
		if w.Job.ID != jobId {
			log.Printf("out of date job of %v, updating to %v", w.Job.ID, jobId)
			w.Job.ID = jobId
		}
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(task api.Task) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		return w.Job.ID != task.Job.ID
	}
}
