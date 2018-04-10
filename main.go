package main

import (
	"github.com/italia/developers-italia-backend/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	cmd.Execute()
}
