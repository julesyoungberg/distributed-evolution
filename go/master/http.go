package master

import (
	"log"
	"net/http"
	"os"

	"github.com/rickyfitts/distributed-evolution/util"

	"golang.org/x/net/websocket"
)

type websocketMessage struct {
	CurrentGeneration uint   `json:"currentGeneration"`
	Output            string `json:"output"`
	TargetImage       string `json:"targetImage"`
}

// TODO handle multiple connections
func (m *Master) subscribe(ws *websocket.Conn) {
	m.ws = ws

	msg := websocketMessage{TargetImage: m.TargetImageBase64}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		log.Println(err)
	}
}

func (m *Master) updateUI(genN uint) {
	util.DPrintf("updating ui")

	m.mu.Lock()
	defer m.mu.Unlock()

	generation, ok := m.Generations[genN]
	if !ok {
		log.Fatalf("error getting generation %v", genN)
	}

	// TODO move this?
	if !generation.Done {
		return
	}

	// get resulting image
	img := generation.Output.Image()

	// send encoded image and current generation
	msg := websocketMessage{
		CurrentGeneration: genN,
		Output:            util.EncodeImage(img),
	}

	if err := websocket.JSON.Send(m.ws, msg); err != nil {
		log.Println(err)
	}
}

// TODO take target image from http
// allow multiple jobs at the same time, take number of workers from request too??
func (m *Master) httpServer() {
	http.Handle("/subscribe", websocket.Handler(m.subscribe))

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
