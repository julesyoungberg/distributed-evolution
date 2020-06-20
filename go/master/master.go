package master

import (
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/util"
	"golang.org/x/net/websocket"
)

// TODO handle multiple jobs
type Master struct {
	Generations       Generations
	NumWorkers        int
	TargetImage       image.Image
	TargetImageBase64 string
	TargetImageHeight int
	TargetImageWidth  int
	Tasks             []api.Task

	mu sync.Mutex
	ws *websocket.Conn
}

// populates the task queue with tasks, where each is a slice of the target image
func (m *Master) generateTasks() {
	util.DPrintf("%v workers available, generating tasks...", m.NumWorkers)

	s := int(math.Floor(math.Sqrt(float64(m.NumWorkers))))

	width, height := util.GetImageDimensions(m.TargetImage)

	m.TargetImageWidth = width
	m.TargetImageHeight = height

	colWidth := int(math.Ceil(float64(width) / float64(s)))
	rowWidth := int(math.Ceil(float64(height) / float64(s)))

	util.DPrintf("splitting image into %v %vpx cols and %v %vpx rows", s, colWidth, s, rowWidth)

	rgbImg := m.TargetImage.(*image.YCbCr)

	m.Tasks = make([]api.Task, m.NumWorkers)

	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			x0 := x * colWidth
			y0 := y * rowWidth

			x1 := int(math.Min(float64(x0+colWidth), float64(width)))
			y1 := int(math.Min(float64(y0+rowWidth), float64(height)))

			rect := image.Rect(x0, y0, x1, y1)

			index := (y * s) + x

			task := api.Task{
				Generation:  1,
				ID:          index,
				Location:    []int{x0, y0},
				Status:      "unstarted",
				TargetImage: util.EncodeImage(rgbImg.SubImage(rect)),
				Type:        "polygons",
			}

			util.DPrintf("creating task with location %v", task.Location)

			m.Tasks[index] = task
		}
	}

	util.DPrintf("%v tasks created", len(m.Tasks))
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

	util.DPrintf("fetching random image...")
	m.TargetImage = util.GetRandomImage()

	util.DPrintf("encoding image...")
	m.TargetImageBase64 = util.EncodeImage(m.TargetImage)

	m.generateTasks()

	go m.httpServer()

	m.rpcServer()
}
