package master

import (
	"github.com/fogleman/gg"

	"github.com/rickyfitts/distributed-evolution/api"
)

type Generation struct {
	Done       bool
	Generation uint
	Output     *gg.Context
	Tasks      []api.Task
}

type Generations = map[uint]Generation

func (m *Master) UpdateGenerations(task api.Task) {
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
			// find and delte the oldest
			var oldest *Generation
			for _, g := range m.Generations {
				if oldest == nil || oldest.Generation > g.Generation {
					oldest = &g
				}
			}
			delete(m.Generations, oldest.Generation)
		}
	}

	// TODO draw to output, then update the ui
	// draw algorithm
	// for each member of the task's population
	//  - cast to a new generic type which implements core drawing functionality
	//  - make sure to offset by the task's Location
	// how to draw a rectangle: (use similar approach for triangle)
	// https://github.com/fogleman/gg/blob/4dc34561c649343936bb2d29e23959bd6d98ab12/context.go#L583
}
