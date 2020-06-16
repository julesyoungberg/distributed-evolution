package master

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type websocketMessage struct {
	TargetImage string `json:"targetImage"`
}

// TODO handle multiple connections
func (m *Master) subscribe(ws *websocket.Conn) {
	msg := websocketMessage{TargetImage: m.TargetImageBase64}

	if err := websocket.JSON.Send(ws, msg); err != nil {
		log.Println(err)
	}

	m.ws = ws
}

func (m *Master) updateUI(generation Generation) {
	// TODO
}

// TODO take target image from http
// allow multiple jobs at the same time, take number of workers from request too??
func (m *Master) HttpServer() {
	http.Handle("/subscribe", websocket.Handler(m.subscribe))

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
