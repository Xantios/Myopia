package main

import (
	"example.com/xantios/tinyproxy/docker"
	"example.com/xantios/tinyproxy/router"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var runningConfig ExportConfig
var logger *logrus.Entry
// var errorPage string

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

	// Setup router
	router.Init()

	// Setup additional hosts
	for _,domain := range runningConfig.domains {
		router.AddHost(domain)
	}

	if runningConfig.debug {
		logger.Debug("Debug mode is enabled")
		router.PrintRouteTable()
	}

	// Subscribe to docker, throw it in a go routine and let it do its thing
	docker.ContainerMapType("HOST","godev.sacredheart.it")
	go docker.Subscribe("")

	logger.Info("Server is starting on "+host)

	http.HandleFunc("/",router.GenericRequestHandler)

	// Wrapped in logger.Fatal in case the listenAndServe call ever fails
	logger.Fatal(http.ListenAndServe(host,nil))
}