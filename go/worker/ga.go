package worker

import (
	"log"
	"os"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (w *Worker) createGA() {
	gaConfig := eaopt.GAConfig{
		NPops:        1,
		PopSize:      w.PopSize,
		HofSize:      1,
		NGenerations: w.NGenerations,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: w.PoolSize,
			},
			MutRate:   w.MutationRate,
			CrossRate: w.CrossRate,
		},
		ParallelEval: false,
	}

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	ga.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	w.ga = ga
}

func createCallback(task api.Task) func(ga *eaopt.GA) {
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

		api.Update(task)
	}
}
