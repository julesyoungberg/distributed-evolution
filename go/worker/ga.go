package worker

import (
	"log"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
)

func createGA() *eaopt.GA {
	ga, err := eaopt.NewDefaultGAConfig().NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}

func createCallback(task api.Task) func(ga *eaopt.GA) {
	return func(ga *eaopt.GA) {
		util.DPrintf("best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)

		task.Generation = ga.Generations
		task.Population = make(eaopt.Individuals, len(ga.Populations[0].Individuals))

		// make a copy of each individual without the context pointer to the worker state
		for i, indv := range ga.Populations[0].Individuals {
			copy := indv
			genome := copy.Genome.Clone()

			if task.Type == "triangles" {
				t := genome.(Triangle)
				t.Context = nil
				genome = t
			}

			copy.Genome = genome

			task.Population[i] = copy
		}

		util.DPrintf("updating master")

		api.Update(task)
	}
}
