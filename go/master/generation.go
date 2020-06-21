package master

import (
	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
	"github.com/rickyfitts/distributed-evolution/go/worker"
)

type Generation struct {
	Done       bool
	Generation uint
	Output     *gg.Context
	Tasks      []api.Task
}

type Generations = map[uint]Generation

func (m *Master) updateGeneration(task *api.Task) Generation {
	util.DPrintf("updating generation %v", task.Generation)

	genN := task.Generation

	generation, ok := m.Generations[genN]

	if ok {
		util.DPrintf("generation %v exists, appending", genN)
		// great, the generation already exists, update it
		generation.Tasks = append(generation.Tasks, *task)
	} else {
		util.DPrintf("generation %v does not exist, creating", genN)
		// this is the first slice of this generation, create it and remove an old one
		generation = Generation{
			Generation: genN,
			Tasks:      []api.Task{*task},
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

	return generation
}

func (m *Master) drawToGeneration(generation Generation, task *api.Task) {
	util.DPrintf("drawing to generation %v with offset %v", generation.Generation, task.Offset)

	s := task.BestFit.Genome.(worker.Shapes)
	s.Draw(generation.Output, task.Offset)

	m.Generations[generation.Generation] = generation
}
