package main

import (
	"github.com/rickyfitts/distributed-evolution/go/api"
	"github.com/rickyfitts/distributed-evolution/go/master"
)

func main() {
	api.Register()
	master.Run()
}
