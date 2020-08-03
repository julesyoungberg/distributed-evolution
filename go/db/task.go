package db

import (
	"fmt"
	"strconv"

	"github.com/rickyfitts/distributed-evolution/go/api"
)

const TaskQueueKey = "task_queue"

func (db *DB) SaveTask(task api.Task) error {
	encoded, err := task.MarshalJSON()
	if err != nil {
		return fmt.Errorf("saving task: %v", err)
	}

	err = db.Set(task.Key(), string(encoded))
	if err != nil {
		return fmt.Errorf("saving task: %v", err)
	}

	return nil
}

func (db *DB) GetTask(id int) (api.Task, error) {
	task := api.Task{ID: id}
	encoded, err := db.Get(task.Key())
	if err != nil {
		return api.Task{}, fmt.Errorf("getting task: %v", err)
	}

	err = task.UnmarshalJSON([]byte(encoded))
	if err != nil {
		return api.Task{}, fmt.Errorf("getting task: %v", err)
	}

	return task, nil
}

func (db *DB) PushTaskID(id int) error {
	_, err := db.Client.RPush(ctx, TaskQueueKey, fmt.Sprint(id)).Result()
	if err != nil {
		return fmt.Errorf("pushing task %v to queue: %v", id, err)
	}

	return nil
}

func (db *DB) PushTask(task api.Task) error {
	err := db.SaveTask(task)
	if err != nil {
		return err
	}

	return db.PushTaskID(task.ID)
}

func (db *DB) PullTask() (api.Task, error) {
	val, err := db.Client.LPop(ctx, TaskQueueKey).Result()
	if err != nil {
		return api.Task{}, fmt.Errorf("pulling task from queue: %v", err)
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return api.Task{}, fmt.Errorf("parsing task id %v: %v", val, err)
	}

	return db.GetTask(id)
}
