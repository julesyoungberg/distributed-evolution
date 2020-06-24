package master

import (
	"image"
	"math"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
	"github.com/rickyfitts/distributed-evolution/go/worker"
)

type Generation struct {
	Done       bool
	Generation uint
	Image      image.Image
	Output     *gg.Context
	Tasks      []api.Task
}

type Generations = map[uint]Generation

// update the generation image with the result of the task
func (m *Master) updateGeneration(task *api.Task) Generation {
	genN := task.Generation

	generation, ok := m.Generations[genN]

	if ok {
		// great, the generation already exists, update it
		generation.Tasks = append(generation.Tasks, *task)
	} else {
		// this is the first slice of this generation, create it and remove an old one
		generation = Generation{
			Generation: genN,
			Tasks:      []api.Task{*task},
		}

		generation.Output = gg.NewContext(m.TargetImage.Width, m.TargetImage.Height)
	}

	if len(generation.Tasks) == len(m.Tasks) {
		// this is the last slice for this generation, mark it as done
		generation.Done = true
	}

	m.Generations[genN] = generation

	return generation
}

// draw the task's output to the generation's drawing context
func (m *Master) drawToGeneration(generation *Generation, task *api.Task) {
	if m.Job.DrawOnce {
		// the generation has already been drawn, draw the image to the output
		img := util.DecodeImage(task.Output)

		centerX := int(math.Round(task.Offset.X + task.Dimensions.X/2.0))
		centerY := int(math.Round(task.Offset.Y + task.Dimensions.Y/2.0))

		generation.Output.DrawImageAnchored(img, centerX, centerY, 0.5, 0.5)
	} else {
		// draw the generation to the output
		task.BestFit.Genome.(worker.Shapes).Draw(generation.Output, task.Offset)
	}
}

func (m *Master) getLatestGeneration() Generation {
	latest := Generation{
		Output: gg.NewContext(m.TargetImage.Width, m.TargetImage.Height),
		Tasks:  m.Tasks,
	}

	for _, t := range m.Tasks {
		if t.Generation > latest.Generation {
			latest.Generation = t.Generation
		}

		m.drawToGeneration(&latest, &t)
	}

	return latest
}
