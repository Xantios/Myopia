package main

import (
	"example.com/xantios/tinyproxy/docker"
	"example.com/xantios/tinyproxy/router"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var runningConfig ExportConfig
var logger *logrus.Entry
var debug bool

func addContainerRoute(hostItem docker.DynamicHost) {

	var mapType router.RouteType
	if strings.HasPrefix(hostItem.Url,"/") {
		mapType = router.MapPath
	} else {
		mapType = router.MapHost
	}

	route := router.Route{
		Name:        "Docker:"+hostItem.ContainerName,
		Source:      hostItem.Url,
		Destination: "http://"+hostItem.Ip+":"+ strconv.Itoa(hostItem.Port),
		MapType:     mapType,
	}

	router.AddRoute(route)
	logger.Warning("Updated route list")
	router.PrintRouteTable()
}

func removeContainerRoute(name string) {
	println("Removing route: "+name)
	router.RemoveRoute("Docker:"+name)
	router.PrintRouteTable()
}

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

	// Set debug flag explicit
	debug = runningConfig.debug

	// Setup router
	router.Init()
	router.SetDebug(debug)

	// Map asset paths
	for _,assets := range runningConfig.assets {

		// URL Sorting in Go is a bit odd. let's amend a slash to the path so the handler knows it's a path
		urlPath := assets.Url+"/"

		fs := http.FileServer(http.Dir(assets.Path))
		handler := http.StripPrefix(urlPath,fs)
		http.Handle(urlPath,handler)
	}

	// Setup additional hosts
	for _,domain := range runningConfig.domains {
		router.AddHost(domain)
	}

	// Print initial route table
	if debug {
		logger.Debug("Debug mode is enabled")
		router.PrintRouteTable()
	}

	// Set container mapping type
	docker.ContainerMapType("HOST","godev.sacredheart.it")

	// Subscribe to docker, convert to route and push to router
	go docker.Subscribe(
		"",
		addContainerRoute,
		removeContainerRoute,
	)

	logger.Info("Server is starting on "+host)

	http.HandleFunc("/",router.GenericRequestHandler)

	// Wrapped in logger.Fatal in case the listenAndServe call ever fails
	logger.Fatal(http.ListenAndServe(host,nil))
}