package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/rickyfitts/distributed-evolution/go/api"
)

const TASK_QUEUE = "queue:task"
const OUTPUT_CHANNEL = "channel:output"

var ctx = context.Background()

type DB struct {
	Client *redis.Client
}

func NewConnection() DB {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_URL"),
		Password: "",
		DB:       0,
	})

	return DB{Client: client}
}

func (db *DB) Set(key string, value string) error {
	err := db.Client.Set(ctx, key, value, 0).Err()
	return err
}

func (db *DB) Get(key string) (string, error) {
	val, err := db.Client.Get(ctx, key).Result()
	return val, err
}

func (db *DB) SaveTask(task api.Task) error {
	encoded, err := task.MarshalJSON()
	if err != nil {
		return err
	}

	err = db.Set(task.Key(), string(encoded))
	if err != nil {
		return fmt.Errorf("saving task: %v", err)
	}

	return nil
}

func (db *DB) GetTask(id uint32) (api.Task, error) {
	log.Printf("getting task %v", id)
	task := api.Task{ID: id}

	json, err := db.Get(task.Key())
	if err != nil {
		e := fmt.Errorf("fetching snapshot for task %v: %v", id, err)
		return task, e
	}

	log.Printf("parsing task %v", id)

	err = task.UnmarshalJSON([]byte(json))

	return task, err
}

func (db *DB) PushTaskID(id uint32) error {
	_, err := db.Client.RPush(ctx, TASK_QUEUE, fmt.Sprint(id)).Result()
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
	val, err := db.Client.LPop(ctx, TASK_QUEUE).Result()
	if err != nil {
		return api.Task{}, fmt.Errorf("pulling task from queue: %v", err)
	}

	log.Printf("from queue: %v", val)

	id, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return api.Task{}, fmt.Errorf("parsing task id %v: %v", val, err)
	}

	return db.GetTask(uint32(id))
}
