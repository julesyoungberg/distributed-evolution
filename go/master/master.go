package master

import (
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

// TODO handle multiple jobs
type Master struct {
	Generations       Generations
	Job               api.Job
	NumWorkers        int
	Outputs           map[int]Generation
	TargetImage       util.Image
	TargetImageBase64 string
	Tasks             []api.Task

	mu         sync.Mutex
	conn       *websocket.Conn
	lastUpdate time.Time
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	log.Printf("%v workers available, generating tasks...", m.NumWorkers)

	s := int(math.Floor(math.Sqrt(float64(m.NumWorkers))))

	width, height := util.GetImageDimensions(m.TargetImage.Image)

	m.TargetImage.Width = width
	m.TargetImage.Height = height

	colWidth := int(math.Ceil(float64(width) / float64(s)))
	rowWidth := int(math.Ceil(float64(height) / float64(s)))

	log.Printf("splitting image into %v %vpx cols and %v %vpx rows", s, colWidth, s, rowWidth)

	m.Tasks = make([]api.Task, m.NumWorkers)

	// create a task for each slice of the image
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			x0 := x * colWidth
			y0 := y * rowWidth
			x1 := int(math.Min(float64(x0+colWidth), float64(width)))
			y1 := int(math.Min(float64(y0+rowWidth), float64(height)))
			rect := image.Rect(x0, y0, x1, y1)

			index := (y * s) + x

			subImg := util.GetSubImage(m.TargetImage.Image, rect)
			bounds := subImg.Bounds()

			task := api.Task{
				Dimensions:  util.Vector{X: float64(bounds.Dx()), Y: float64(bounds.Dy())},
				Generation:  1,
				ID:          index,
				Offset:      util.Vector{X: float64(x0), Y: float64(y0)},
				Status:      "unstarted",
				TargetImage: util.EncodeImage(subImg),
				Type:        "polygons",
			}

			m.Tasks[index] = task
		}
	}

	log.Printf("%v tasks created", len(m.Tasks))
}

func Run() {
	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		log.Fatal("error parsing NUM_WORKERS: ", err)
	}

	jobId := uuid.New().ID()

	m := Master{
		Generations: Generations{},
		NumWorkers:  numWorkers,
		Outputs:     map[int]Generation{},
		Job: api.Job{
			ID:           jobId,
			DrawOnce:     true,
			CrossRate:    0.2,
			MutationRate: 0.021,
			NumShapes:    200,
			OutputMode:   "combined",
			OverDraw:     10,
			PoolSize:     10,
			PopSize:      50,
			ShapeSize:    20,
		},
		TargetImage: util.Image{},
	}

	log.Print("fetching random image...")
	m.TargetImage.Image = util.GetRandomImage()

	log.Print("encoding image...")
	m.TargetImageBase64 = util.EncodeImage(m.TargetImage.Image)

	m.generateTasks()

	go m.httpServer()

	m.rpcServer()
}
