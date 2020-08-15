package worker

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

func saveImageToFile(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

func saveBase64ToPng(base64, filename string) error {
	img, err := util.DecodeImage(base64)
	if err != nil {
		return err
	}

	return saveImageToFile(img, filename)
}

func savePreviousOutput(output string) error {
	return saveBase64ToPng(output, "previous_output.png")
}

func saveTargetImage(target string) error {
	return saveBase64ToPng(target, "target.png")
}

func saveCurrentBest(task api.Task) error {
	context, err := getTaskContext(task, []color.RGBA{{0, 0, 0, 0}})
	if err != nil {
		return fmt.Errorf("error saving current best: getting worker task: %v", err)
	}

	for i, member := range task.Population {
		s := member.Genome.(api.Shapes)
		s.Context = context
		if _, err := s.Evaluate(); err != nil {
			log.Printf("error evaluating member %v", i)
		}
	}

	return saveImageToFile(context.BestFit.Output, "current_best.png")
}

func createTestingCallback(ctx *api.TaskContext) func(ga *eaopt.GA) {
	return func(ga *eaopt.GA) {
		task, err := ctx.EnrichTask(ga)
		if err != nil {
			log.Printf("error enriching task: %v", err)
			return
		}

		ctx.Task = task
	}
}

func createTestingEarlyStop(generations uint) func(ga *eaopt.GA) bool {
	return func(ga *eaopt.GA) bool {
		return ga.Generations >= generations
	}
}

func run(t *testing.T, task api.Task, generations uint) *api.TaskContext {
	task.Job.NumGenerations = generations
	context, err := getTaskContext(task, []color.RGBA{})
	if err != nil {
		t.Fatal(err)
	}

	ga := CreateGA(task.Job)

	ga.Callback = createTestingCallback(context)
	ga.EarlyStop = createTestingEarlyStop(generations)
	factory := api.GetShapesFactory(context, task.Population)

	err = ga.Minimize(factory)
	if err != nil {
		t.Fatal(err)
	}

	return context
}

func fatalError(t *testing.T, f func() error) {
	if err := f(); err != nil {
		t.Fatal(err)
	}
}

func logError(t *testing.T, f func() error) {
	if err := f(); err != nil {
		t.Log(err)
	}
}

func TestRecover(t *testing.T) {
	file, err := os.Open("snapshot.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	var task api.Task
	fatalError(t, func() error { return task.UnmarshalJSON(bytes) })

	logError(t, func() error { return saveTargetImage(task.TargetImage) })
	// logError(t, func() error { return savePreviousOutput(task.Output) })
	// logError(t, func() error { return saveCurrentBest(task) })

	// ctx := run(t, task, 1)
	// logError(t, func() error { return saveImageToFile(ctx.BestFit.Output, "output.png") })

	ctx := run(t, task, 500)
	logError(t, func() error { return saveImageToFile(ctx.BestFit.Output, "initial_output.png") })

	bytes, err = ctx.Task.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var task2 api.Task
	fatalError(t, func() error { return task2.UnmarshalJSON(bytes) })

	logError(t, func() error { return savePreviousOutput(task2.Output) })
	logError(t, func() error { return saveCurrentBest(task2) })

	ctx = run(t, task2, 1)
	logError(t, func() error { return saveImageToFile(ctx.BestFit.Output, "output.png") })
}
