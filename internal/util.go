package internal

import (
	"log"
	"time"
)

func Timed(actionName string, f func()) {
	start := time.Now()
	f()
	log.Printf("Completed %s in %v", actionName, time.Since(start))
}
