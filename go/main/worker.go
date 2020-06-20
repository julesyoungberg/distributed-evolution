package main

import "github.com/rickyfitts/distributed-evolution/go/worker"

func main() {
	worker.Register()
	worker.Run()
}
