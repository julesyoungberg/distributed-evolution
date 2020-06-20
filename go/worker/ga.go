package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

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

func (w *Worker) createCallback(task api.Task) func(ga *eaopt.GA) {
	return func(ga *eaopt.GA) {
		util.DPrintf("best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)

		bestFit := ga.HallOfFame[0]
		var genome eaopt.Genome

		t := bestFit.Genome.(Shapes)
		genome = t.CloneWithoutContext()

		// for _, m := range t.Members {
		// 	util.DPrintf("vertices: %v", m.Vertices)
		// 	util.DPrintf("color: %v", m.Color)
		// }

		bestFit.Genome = genome

		task.BestFit = bestFit
		task.Generation = ga.Generations

		util.DPrintf("updating master")

		jobId, err := api.Update(task)
		if err != nil {
			log.Printf("error %v", err)
		}

		if w.Job.ID != jobId {
			log.Printf("out of date job of %v, updating to %v", w.Job.ID, jobId)
			w.Job.ID = jobId
		}
	}
}

func (w *Worker) createEarlyStop(task api.Task) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		if w.Job.ID != task.Job.ID {
			log.Printf("earlystop")
		}

		return w.Job.ID != task.Job.ID
	}
}
