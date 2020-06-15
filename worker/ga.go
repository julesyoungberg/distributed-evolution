package main

import (
	"log"

	"github.com/MaxHalford/eaopt"
)

func createGA() *eaopt.GA {
	gaConfig := eaopt.NewDefaultGAConfig()

	gaConfig.NPops = 1
	gaConfig.NGenerations = 1

	ga, err := gaConfig.NewGA()
	if err != nil {
		log.Fatal("error creating ga: ", err)
	}

	return ga
}
