package master

import (
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/rickyfitts/distributed-evolution/api"
	"github.com/rickyfitts/distributed-evolution/util"
	"golang.org/x/net/websocket"
)

// TODO handle multiple jobs
type Master struct {
	Generations       Generations
	NumWorkers        int
	TargetImage       image.Image
	TargetImageBase64 string
	Tasks             []api.Task

	mu sync.Mutex
	ws *websocket.Conn
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) GenerateTasks() {
	s := math.Floor(math.Sqrt(float64(m.NumWorkers)))

	width, height := util.GetImageDimensions(m.TargetImage)

	cols := int(math.Ceil(float64(width) / s))
	rows := int(math.Ceil(float64(height) / s))

	rgbImg := m.TargetImage.(*image.YCbCr)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			x0 := x * int(s)
			y0 := y * int(s)

			x1 := int(math.Min(float64(x0)+s, float64(width)))
			y1 := int(math.Min(float64(y0)+s, float64(height)))

			rect := image.Rect(x0, y0, x1, y1)

			task := api.Task{
				Generation:  1,
				ID:          (y * cols) + x,
				Location:    []int{x0, y0},
				Status:      "unstarted",
				TargetImage: util.EncodeImage(rgbImg.SubImage(rect)),
				Type:        "triangles",
			}

			m.Tasks = append(m.Tasks, task)
		}
	}
}

func Run() {
	numWorkers, err := strconv.Atoi(os.Getenv("NUM_WORKERS"))
	if err != nil {
		log.Fatal("error parsing NUM_WORKERS: ", err)
	}

	m := Master{
		Generations: make(Generations, 3),
		NumWorkers:  numWorkers,
	}

	m.TargetImage = util.GetRandomImage()
	m.TargetImageBase64 = util.EncodeImage(m.TargetImage)

	m.GenerateTasks()

	go m.HttpServer()

	m.RpcServer()
}
