package api

import (
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Output struct {
	Fitness float64
	Output  image.Image
}

type TaskContext struct {
	BestFit     Output
	Edges       image.Image
	GenOffset   uint
	Mu          sync.Mutex
	Palette     []color.RGBA
	TargetImage util.Image
	Task        Task
}

type Job struct {
	Complete       bool      `json:"complete"`
	CompletedAt    time.Time `json:"completedAt"`
	CrossRate      float64   `json:"crossRate"`
	DetectEdges    bool      `json:"detectEdges"`
	ID             int       `json:"ID"`
	MutationRate   float64   `json:"mutationRate"`
	NumColors      int       `json:"numColors"`
	NumGenerations uint      `json:"numGenerations"`
	NumShapes      int       `json:"numShapes"`
	OverDraw       int       `json:"overDraw"`
	PaletteType    string    `json:"paletteType"`
	PoolSize       uint      `json:"poolSize"`
	PopSize        uint      `json:"popSize"`
	Quantization   int       `json:"quantization"`
	ShapeSize      uint      `json:"shapeSize"`
	ShapeType      string    `json:"shapeType"`
	ShapesPerSlice int       `json:"shapesPerSlice"`
	StartedAt      time.Time `json:"startedAt"`
	TargetImage    string    `json:"targetImage"`
}

type Task struct {
	BestFit            eaopt.Individual  `json:"-"`
	Dimensions         util.Vector       `json:"dimensions"`
	Edges              string            `json:"edges"`
	Generation         uint              `json:"generation"`
	ID                 int               `json:"ID"`
	Job                Job               `json:"job"`
	Output             string            `json:"output"`
	Population         eaopt.Individuals `json:"-"`
	Position           util.Vector       `json:"position"`
	ScaledQuantization int               `json:"quantization"`
	ShapeType          string            `json:"shapeType"`
	TargetImage        string            `json:"targetImage"`
}

type TaskState struct {
	Attempt     int       `json:"attempt"`
	Complete    bool      `json:"complete"`
	CompletedAt time.Time `json:"completedAt"`
	Fitness     float64   `json:"fitness"`
	Generation  uint      `json:"generation"`
	ID          int       `json:"ID"`
	JobID       int       `json:"jobID"`
	LastUpdate  time.Time `json:"lastUpdate"`
	StartedAt   time.Time `json:"startedAt"`
	Status      string    `json:"status"`
	Thread      int       `json:"thread"`
	WorkerID    uint32    `json:"workerID"`
}

type MasterSnapshot struct {
	Job         Job               `json:"job"`
	TargetImage string            `json:"targetImage"`
	Tasks       map[int]TaskState `json:"tasks"`
}

func Register() {
	gob.Register(color.RGBA{})
	gob.Register(image.YCbCr{})

	gob.Register(Circle{})
	gob.Register(Polygon{})
	gob.Register(Triangle{})
	gob.Register(Shapes{})
}

func (t Task) Key() string {
	return fmt.Sprintf("task:%v", t.ID)
}

func (t Task) UpdateMaster(worker uint32, thread int, status string) error {
	fitness := t.BestFit.Fitness
	if fitness != 0 {
		fitness = 1 / fitness
	}

	return Update(TaskState{
		Fitness:    fitness,
		Generation: t.Generation,
		ID:         t.ID,
		JobID:      t.Job.ID,
		Status:     status,
		Thread:     thread,
		WorkerID:   worker,
	})
}

func (ctx *TaskContext) EnrichTask(ga *eaopt.GA) (Task, error) {
	output := ctx.BestFit.Output
	if output == nil {
		return Task{}, fmt.Errorf("best fit output is nil")
	}

	encoded, err := util.EncodeImage(output)
	if err != nil {
		return Task{}, fmt.Errorf("encoding output: %v", err)
	}

	bestFit := ga.HallOfFame[0]
	bestFit.Genome = Shapes{}

	task := ctx.Task
	task.Output = encoded
	task.BestFit = bestFit

	task.Population = ga.Populations[0].Individuals

	return task, nil
}
