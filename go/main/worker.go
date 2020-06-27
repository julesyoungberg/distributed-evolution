package main

import (
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/worker"
)

func main() {
	api.Register()
	worker.Run()
}
