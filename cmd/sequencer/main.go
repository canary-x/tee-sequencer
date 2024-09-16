package main

import (
	"log"

	"github.com/canary-x/tee-sequencer/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		log.Fatalf("Fatal error: %+v\n", err)
	}
}
