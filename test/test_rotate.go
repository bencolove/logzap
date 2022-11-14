//go:build rotate

package main

import (
	log "github.com/bencolove/logzap"
)

// go run -tags rotage test/test_rotate.go
func main() {
	config := &log.RotateOption{
		Filepath: "./rlog.log",
		// 1 M
		MaxSize: 1,
		// 1 day
		MaxAge:    1,
		MaxBackup: 10,
		Compress:  true,
	}

	logger := log.New(config.Writer(), log.InfoLevel)
	log.ResetDefault(logger)
	defer log.Sync()

	log.Info("demo_rotatefile", log.String("app", "started"), log.Int("version", 2))
}
