package api

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"time"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Job struct {
	CrossRate    float64 `json:"crossRate"`
	ID           uint32  `json:"ID"`
	MutationRate float64 `json:"mutationRate"`
	NumShapes    uint    `json:"numShapes"`
	OverDraw     int     `json:"overDraw"`
	PoolSize     uint    `json:"poolSize"`
	PopSize      uint    `json:"popSize"`
	ShapeSize    uint    `json:"shapeSize"`
	TargetImage  string  `json:"targetImage"`
}

type GetTaskArgs struct {
	WorkerID uint32
}

type Task struct {
	BestFit     eaopt.Individual `json:"bestFit"`
	Dimensions  util.Vector      `json:"dimensions"`
	Generation  uint             `json:"generation"`
	ID          uint32           `json:"ID"`
	Job         Job              `json:"-"`
	LastUpdate  time.Time        `json:"lastUpdate"`
	Linked      []int            `json:"linked"`
	Offset      util.Vector      `json:"offset"`
	Output      string           `json:"output"`
	Population  eaopt.Population `json:"population"`
	Status      string           `json:"status"`
	TargetImage string           `json:"-"`
	Thread      int              `json:"thread"`
	Type        string           `json:"type"`
	WorkerID    uint32           `json:"workerID"`
}

func Update(args Task) (Task, error) {
	var reply Task

	err := Call("Master.Update", &args, &reply)

	return reply, err
}

// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func Call(rpcname string, args interface{}, reply interface{}) error {
	c, err := rpc.DialHTTP("tcp", os.Getenv("MASTER_URL"))
	if err != nil {
		return err
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err != nil {
		return err
	}

	return nil
}

func (t Task) Key() string {
	return fmt.Sprintf("task:%v", t.ID)
}

func (t Task) ToJson() (string, error) {
	encoded, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("error encoding task %v: %v", t.ID, err)
	}

	return string(encoded), nil
}

func ParseTaskJson(s string) (Task, error) {
	bytes := []byte(s)
	var task Task

	err := json.Unmarshal(bytes, &task)
	if err != nil {
		return Task{}, fmt.Errorf("error parsing task: %v", err)
	}

	return task, nil
}

func (t Task) UpdateMaster(status string) (Task, error) {
	return Update(Task{
		ID:         t.ID,
		Generation: t.Generation,
		Job:        Job{ID: t.Job.ID},
		Status:     status,
		Thread:     t.Thread,
		WorkerID:   t.ID,
	})
}
