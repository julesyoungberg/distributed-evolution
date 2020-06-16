package master

import (
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type websocketMessage struct {
	TargetImage string `json:"targetImage"`
}

func (m *Master) subscribe(ws *websocket.Conn) {
	for {
		msg := websocketMessage{TargetImage: m.targetImageBase64}

		if err := websocket.JSON.Send(ws, msg); err != nil {
			log.Println(err)
			break
		}

		time.Sleep(10 * time.Second)
	}
}

func (m *Master) httpServer() {
	http.Handle("/subscribe", websocket.Handler(m.subscribe))

	port := os.Getenv("HTTP_PORT")

	log.Printf("listening for HTTP on port %v\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
