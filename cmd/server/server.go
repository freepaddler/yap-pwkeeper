package main

import "yap-pwkeeper/internal/pkg/logger"

var (
	// go build -ldflags " \
	// -X 'main.buildVersion=$(git describe --tag 2>/dev/null)' \
	// -X 'main.buildDate=$(date)' \
	// -X 'main.buildCommit=$(git rev-parse --short HEAD)' \
	// "
	buildVersion, buildDate, buildCommit = "N/A", "N/A", "N/A"
)

func main() {

	logger.SetMode(logger.ModeProd)
	logger.Log().Info("Server Started")
	logger.Log().Info("Server Stopped")
}
