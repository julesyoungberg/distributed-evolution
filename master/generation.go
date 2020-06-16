package master

import (
	"log"

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

func (m *Master) UpdateGenerations(task api.Task) uint {
	m.mu.Lock()
	defer m.mu.Unlock()

	genN := task.Population.Generations

	generation, ok := m.Generations[genN]

	if ok {
		// great, the generation already exists, update it
		generation.Tasks = append(generation.Tasks, task)

		if len(generation.Tasks) == len(m.Tasks) {
			// this is the last slice for this generation, mark it as done
			generation.Done = true
		}
	} else {
		// this is the first slice of this generation, create it and remove an old one
		generation = Generation{
			Generation: genN,
			Tasks:      []api.Task{task},
		}

		generation.Output = gg.NewContext(1000, 1000)

		m.Generations[genN] = generation

		if len(m.Generations) > 3 {
			// find and delete the oldest
			var oldest *Generation
			for _, g := range m.Generations {
				if oldest == nil || oldest.Generation > g.Generation {
					oldest = &g
				}
			}
			delete(m.Generations, oldest.Generation)
		}
	}

	return genN
}

func (m *Master) DrawToGeneration(genN uint, task api.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	for _, member := range task.Population.Individuals {
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
}
