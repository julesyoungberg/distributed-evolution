package master

import (
	"log"
	"net/http"
	"os"
	"time"

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
	m.mu.Lock()
	defer m.mu.Unlock()

	util.DPrintf("new websocket connection request")

	m.ws = ws

	msg := websocketMessage{TargetImage: m.TargetImageBase64}

	util.DPrintf("sending data")

	if err := websocket.JSON.Send(ws, msg); err != nil {
		log.Println(err)
	}

	for {
		// keep the connection open
		m.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
		m.mu.Lock()
		// TODO send keep alive and check if the connection has been closed
	}
}

func (m *Master) updateUI(genN uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	util.DPrintf("updating ui with generation %v", genN)

	if m.ws == nil {
		util.DPrintf("no open ui connections")
		return
	}

	generation, ok := m.Generations[genN]
	if !ok {
		log.Fatalf("error getting generation %v", genN)
	}

	// TODO move this?
	if !generation.Done {
		util.DPrintf("error, generation not done!")
		return
	}

	util.DPrintf("encoding output image for generation %v", genN)

	// get resulting image
	img := generation.Output.Image()

	// send encoded image and current generation
	msg := websocketMessage{
		CurrentGeneration: genN,
		Output:            util.EncodeImage(img),
	}

	util.DPrintf("sending generation %v update to ui", genN)

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
