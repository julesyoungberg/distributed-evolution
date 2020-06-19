package master

import (
	"log"

	"github.com/rickyfitts/distributed-evolution/util"

	"github.com/fogleman/gg"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/worker"
)

type Generation struct {
	Done       bool
	Generation uint
	Output     *gg.Context
	Tasks      []api.Task
}

type Generations = map[uint]Generation

func (m *Master) updateGenerations(task api.Task) uint {
	util.DPrintf("updating generation %v", task.Generation)

	genN := task.Generation

	generation, ok := m.Generations[genN]

	if ok {
		util.DPrintf("generation %v exists, appending", genN)
		// great, the generation already exists, update it
		generation.Tasks = append(generation.Tasks, task)
	} else {
		util.DPrintf("generation %v does not exist, creating", genN)
		// this is the first slice of this generation, create it and remove an old one
		generation = Generation{
			Generation: genN,
			Tasks:      []api.Task{task},
		}

		generation.Output = gg.NewContext(m.TargetImageWidth, m.TargetImageHeight)

		// if len(m.Generations) > 50 {
		// 	delete(m.Generations, genN-50)
		// }
	}

	util.DPrintf("generation %v recieved %v out of %v tasks", genN, len(generation.Tasks), len(m.Tasks))

	if len(generation.Tasks) == len(m.Tasks) {
		util.DPrintf("all tasks complete, marking generation %v as done", genN)
		// this is the last slice for this generation, mark it as done
		generation.Done = true
	}

	m.Generations[genN] = generation

	return genN
}

func (m *Master) drawToGeneration(genN uint, task api.Task) {
	util.DPrintf("drawing to generation %v", genN)

	generation, ok := m.Generations[genN]
	if !ok {
		// wtf
		log.Fatalf("error getting generation %v", genN)
	}

	offset := util.Vector{X: float64(task.Location[0]), Y: float64(task.Location[1])}

	util.DPrintf("drawing with offset %v", offset)

	if task.Type == "triangles" {
		t := task.BestFit.Genome.(worker.Shapes)
		t.Draw(generation.Output, offset)
	}

	m.Generations[genN] = generation
}
