package main

import log "github.com/sirupsen/logrus"

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	log.Println("starting...")
}
