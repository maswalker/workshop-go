package main

import (
	"go-seven/pkg/log"
)

func main() {
	log.Info("demo3:", log.String("app", "start ok"),
		log.Int("major version", 3))
	log.Error("demo3:", log.String("app", "crash"),
		log.Int("reason", -1))
}
