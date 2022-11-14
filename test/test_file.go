//go:build file

package main

import (
	"os"

	log "github.com/bencolove/logzap"
)

// go run -tags file test/test_file.go
func main() {
	file, err := os.OpenFile("./demo.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, log.InfoLevel)
	log.ResetDefault(logger)
	defer log.Sync()

	log.Info("demo_writefile", log.String("app", "started"), log.Int("version", 2))
}
