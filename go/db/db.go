package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

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
	encoded, err := task.ToJson()
	if err != nil {
		return err
	}

	err = db.Set(task.Key(), encoded)
	if err != nil {
		return fmt.Errorf("error saving task: ", err)
	}

	return nil
}

func (db *DB) GetTask(id uint32) (api.Task, error) {
	task := api.Task{ID: id}

	json, err := db.Get(task.Key())
	if err != nil {
		e := fmt.Errorf("error fetching snapshot for task %v: %v", id, err)
		return api.Task{}, e
	}

	return api.ParseTaskJson(json)
}

func (db *DB) PushTaskID(id uint32) error {
	_, err := db.Client.RPush(ctx, TASK_QUEUE, fmt.Sprint(id)).Result()
	if err != nil {
		return fmt.Errorf("error pushing task %v to queue: %v", id, err)
	}

	return nil
}

func (db *DB) PushTask(task api.Task) error {
	db.SaveTask(task)
	return db.PushTaskID(task.ID)
}

func (db *DB) PullTask() (api.Task, error) {
	val, err := db.Client.BLPop(ctx, 200*time.Millisecond, TASK_QUEUE).Result()
	if err != nil {
		return api.Task{}, fmt.Errorf("error pulling task from queue: %v", err)
	}

	id, err := strconv.ParseUint(val[0], 10, 32)
	if err != nil {
		return api.Task{}, fmt.Errorf("error parsing task id %v: %v", val, err)
	}

	return db.GetTask(uint32(id))
}
