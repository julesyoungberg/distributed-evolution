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
		PopSize:      w.Task.Job.PopSize,
		HofSize:      1,
		NGenerations: w.NGenerations,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: w.Task.Job.PoolSize,
			},
			MutRate:   w.Task.Job.MutationRate,
			CrossRate: w.Task.Job.CrossRate,
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
func (w *Worker) createCallback() func(ga *eaopt.GA) {
	// send the currrent best fit and other data to the master
	return func(ga *eaopt.GA) {
		w.SaveTaskSnapshot()

		// get best fit
		bestFit := ga.HallOfFame[0]

		// if draw once is active add the output and clear the genome
		if w.Task.Job.DrawOnce {
			w.mu.Lock()
			output := w.BestFit.Output
			w.mu.Unlock()

			if output == nil {
				log.Printf("error! best fit image is nil at generation %v - bestFit: %v", ga.Generations, w.BestFit)
				return
			}

			w.Task.Output = util.EncodeImage(output)
			bestFit.Genome = Shapes{}

			// clear state
			w.BestFit = Output{}
		} else {
			// otherwise just add the genome
			bestFit.Genome = bestFit.Genome.(Shapes).CloneForSending()
		}

		// add data to the task
		w.Task.BestFit = bestFit
		w.Task.Generation = ga.Generations

		// send results to master
		jobId, err := api.Update(w.Task)
		if err != nil {
			log.Print("error ", err)
		}

		// if the master responded with a different job id we are out of date
		if w.Task.Job.ID != jobId {
			log.Printf("out of date job of %v, updating to %v", w.Task.Job.ID, jobId)
			w.Task.Job.ID = jobId
			return
		}

		// TODO: hamdle newly assigned / linked tasks
	}
}

// returns a closure to check if the algorithm should stop (ie the job has changed)
func (w *Worker) createEarlyStop(id uint32) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		return w.Task.Job.ID != id
	}
}
