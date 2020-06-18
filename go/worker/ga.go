package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func createGA() *eaopt.GA {
	gaConfig := eaopt.GAConfig{
		NPops:        1,
		PopSize:      100,
		HofSize:      1,
		NGenerations: 50,
		Model: eaopt.ModGenerational{
			Selector: eaopt.SelTournament{
				NContestants: 3,
			},
			MutRate:   0.5,
			CrossRate: 0.7,
		},
		ParallelEval: false,
	}

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}

func createCallback(task api.Task) func(ga *eaopt.GA) {
	return func(ga *eaopt.GA) {
		util.DPrintf("best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)

		bestFit := ga.HallOfFame[0]
		var genome eaopt.Genome

		if task.Type == "triangles" {
			t := bestFit.Genome.(Triangles)
			genome = t.CloneWithoutContext()
		}

		bestFit.Genome = genome

		task.BestFit = bestFit
		task.Generation = ga.Generations

		util.DPrintf("updating master")

		api.Update(task)
	}
}
