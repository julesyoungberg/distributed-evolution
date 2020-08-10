package master

import (
	"image"
	"log"
	"time"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

func (m *Master) allStale() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	log.Printf("checking for staleness")

	for _, t := range m.Tasks {
		if t.Status != "stale" {
			log.Printf("task %v is %v", t.ID, t.Status)
			return false
		}
	}

	log.Printf("all tasks are stale")
	return true
}

func (m *Master) setTargetImage(image image.Image) {
	width, height := util.GetImageDimensions(image)
	m.TargetImage = util.Image{
		Image:  image,
		Height: height,
		Width:  width,
	}
}

func (m *Master) startRandomJob() {
	log.Print("fetching random image...")
	image := util.GetRandomImage()

	log.Print("encoding image...")
	encodedImg, err := util.EncodeImage(image)
	if err != nil {
		log.Fatal(err)
	}

	m.TargetImageBase64 = encodedImg
	m.setTargetImage(image)
	m.Job.StartedAt = time.Now()

	go m.generateTasks()
}

func (m *Master) stopJob() {
	log.Printf("transitioning")

	m.mu.Lock()
	m.transitioning = true
	m.mu.Unlock()

	log.Printf("clearing the task queue")

	// clear the queue of all tasks
	for {
		_, err := m.db.PullTask()
		if err != nil {
			break
		}
	}

	m.mu.Lock()

	// mark any qeued tasks as stale
	for _, task := range m.Tasks {
		if task.Status == "queued" {
			task.Status = "stale"
		}
	}

	m.mu.Unlock()

	log.Printf("waiting for workers to stop")

	// wait for all the workers to stop
	for !m.allStale() {
		time.Sleep(time.Second)
	}

	m.mu.Lock()
	m.transitioning = false
	m.mu.Unlock()

	log.Printf("transition complete")
}
