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
	m.mu.Lock()
	defer m.mu.Unlock()

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

		generation.Output = gg.NewContext(1000, 1000)

		if len(m.Generations) > 3 {
			delete(m.Generations, genN-3)
		}
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
	m.mu.Lock()
	defer m.mu.Unlock()

	util.DPrintf("drawing to generation %v", genN)

	generation, ok := m.Generations[genN]
	if !ok {
		// wtf
		log.Fatalf("error getting generation %v", genN)
	}

	// get offset from task
	offsetX := float64(task.Location[0])
	offsetY := float64(task.Location[1])

	dc := generation.Output

	// draw each member of the population
	for _, member := range task.Population {
		// check the type of task
		if task.Type == "triangles" {
			t := member.Genome.(worker.Triangle)

			// draw triangle
			dc.NewSubPath()
			dc.MoveTo(t.Vertices[0][0]+offsetX, t.Vertices[0][1]+offsetY)
			dc.LineTo(t.Vertices[1][0]+offsetX, t.Vertices[1][1]+offsetY)
			dc.LineTo(t.Vertices[2][0]+offsetX, t.Vertices[2][1]+offsetY)
			dc.ClosePath()

			dc.SetColor(t.Color)
			dc.Fill()
		}
		// TODO implement more shapes
	}

	m.Generations[genN] = generation
}