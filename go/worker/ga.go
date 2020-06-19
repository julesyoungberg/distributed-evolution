package worker

import (
	"log"
	"os"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func createGA(crossRate, mutationRate float64) *eaopt.GA {
	gaConfig := eaopt.GAConfig{
		NPops:        1,
		PopSize:      100,
		HofSize:      1,
		NGenerations: 1000,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: 20,
			},
			MutRate:   mutationRate,
			CrossRate: crossRate,
		},
		ParallelEval: false,
	}

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	ga.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return ga
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
