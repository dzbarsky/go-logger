package main

import (
	"time"

	loglib "logger/logger"
)

func main() {
	logger := loglib.New()
	logger.Tag("version", "deadbeef")

	loglib.SetSeverity("", loglib.Debug)

	go toggleSeverities()

	for {
		time.Sleep(300 * time.Millisecond)

		logger.Info("root info")
		logger.Debug("root debug")
		logger.Error("root error")
		serverCode(logger)
	}
}

func serverCode(logger *loglib.Logger) {
	logger = logger.PushScope("server")
	logger.Tag("service", "kv-store")

	logger.Info("server info")
	logger.Debug("server debug")
	logger.Error("server error")

	logger = logger.PushScope("server.init_block")
	logger.Info("server init info")
	logger.Debug("server init debug")
	logger.Error("server init error")
	logger = logger.PopScope()

	logger.Info("server info")
	logger.Debug("server debug")
	logger.Error("server error")
}

func toggleSeverities() {
	for {
		for _, s := range []loglib.Severity{loglib.Debug, loglib.Error} {
			loglib.SetSeverity("server", s)
			time.Sleep(1 * time.Second)
		}
	}
}


