package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type websocketMessage struct {
	TargetImage string `json:"targetImage"`
}

func (m *Master) socket(ws *websocket.Conn) {
	for {
		m := websocketMessage{TargetImage: m.targetImageBase64}

		if err := websocket.JSON.Send(ws, m); err != nil {
			log.Println(err)
			break
		}

		time.Sleep(10 * time.Second)
	}
}

func (m *Master) httpServer() {
	http.Handle("/subscribe", websocket.Handler(m.socket))

	port := os.Getenv("HTTP_PORT")

	fmt.Printf("http listening on port %v\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
