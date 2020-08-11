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

type WorkerTask struct {
	BestFit     Output
	Edges       image.Image
	GenOffset   uint
	Mu          sync.Mutex
	Palette     []color.RGBA
	TargetImage util.Image
	Task        Task
}

type Job struct {
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
	ID         int       `json:"ID"`
	Generation uint      `json:"generation"`
	JobID      int       `json:"jobID"`
	LastUpdate time.Time `json:"lastUpdate"`
	Status     string    `json:"status"`
	Thread     int       `json:"thread"`
	WorkerID   uint32    `json:"workerID"`
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
	return Update(TaskState{
		ID:         t.ID,
		Generation: t.Generation,
		JobID:      t.Job.ID,
		Status:     status,
		Thread:     thread,
		WorkerID:   worker,
	})
}
