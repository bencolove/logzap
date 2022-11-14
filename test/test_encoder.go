//go:builder encoder

package main

import (
	"os"

	log "github.com/bencolove/logzap"
)

// go run -tags encoder test/test_encoder.go
func main() {
	logger := log.New(os.Stdout, log.InfoLevel,
		log.WithCaller(true),
		log.WithTimeFormat("2006-01-02 15:04:05"),
	)
	log.ResetDefault(logger)
	defer log.Sync()

	log.Info("demo_encoder", log.String("app", "started"), log.Int("version", 2))
}
