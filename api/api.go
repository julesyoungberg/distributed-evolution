package api

import (
	"image"
	"time"
)

type EmptyMessage struct{}

type Task struct {
	Generation  int
	ID          int
	Location    image.Rectangle
	Started     time.Time
	TargetImage string
}
