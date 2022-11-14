//go:build tee

// test writing trivial logs(less critical then INFO) to access.log and
// writing more severe logs(more critical then INFO) to error.log
package main

import (
	"os"

	log "github.com/bencolove/logzap"
)

// go run -tags tee test/test_tee.go
func main() {
	file1, err := os.OpenFile("./access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	file2, err := os.OpenFile("./error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var tops = []log.TeeOption{
		{
			W: file1,
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
		{
			W: file2,
			Lef: func(lvl log.Level) bool {
				return lvl > log.InfoLevel
			},
		},
	}

	logger := log.NewTee(tops)
	log.ResetDefault(logger)

	log.Info("demo3:", log.String("app", "start ok"),
		log.Int("major version", 3))
	log.Error("demo3:", log.String("app", "crash"),
		log.Int("reason", -1))

}
