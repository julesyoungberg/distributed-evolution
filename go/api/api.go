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
	GenOffset   uint
	Mu          sync.Mutex
	TargetImage util.Image
	Task        Task
}

type Job struct {
	CrossRate      float64 `json:"crossRate"`
	ID             int     `json:"ID"`
	MutationRate   float64 `json:"mutationRate"`
	NumShapes      int     `json:"numShapes"`
	OverDraw       int     `json:"overDraw"`
	PoolSize       uint    `json:"poolSize"`
	PopSize        uint    `json:"popSize"`
	ShapeSize      uint    `json:"shapeSize"`
	ShapesPerSlice int     `json:"shapesPerSlice"`
	TargetImage    string  `json:"targetImage"`
}

type Task struct {
	BestFit     eaopt.Individual  `json:"-"`
	Dimensions  util.Vector       `json:"dimensions"`
	Generation  uint              `json:"generation"`
	ID          int               `json:"ID"`
	Job         Job               `json:"job"`
	Offset      util.Vector       `json:"offset"`
	Output      string            `json:"output"`
	Population  eaopt.Individuals `json:"-"`
	TargetImage string            `json:"targetImage"`
	Type        string            `json:"type"`
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
