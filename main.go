package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var runningConfig ExportConfig
var logger *logrus.Entry
var errorPage string

func main() {

	// Pull config
	runningConfig = GetConf("config.yaml")
	host := env("host","0.0.0.0")+":"+env("port","42069")

	// Setup logging
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logger = logrus.WithFields(logrus.Fields{
		"service": "proxy",
	})

	if runningConfig.debug {
		logger.Debug("Debug mode is enabled")
	}

	logger.Info("Server is starting on "+host)

	// Preload pages
	errorPage = LoadInternalAsset("./assets/error.page.html")

	http.HandleFunc("/",GenericRequestHandler)

	// Wrapped in logger.Fatal in case the listenAndServe call ever fails
	logger.Fatal(http.ListenAndServe(host,nil))
}