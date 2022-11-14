//go:build run

package main

import (
	log "github.com/bencolove/logzap"
)

// go run -tags run test/test_package.go
func main() {
	defer log.Sync()

	log.Info("demo1", log.String("app", "started"), log.Int("majorVersion", 2))

}
