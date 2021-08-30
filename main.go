package main

import (
	"crypto/tls"
	"example.com/xantios/tinyproxy/docker"
	"example.com/xantios/tinyproxy/router"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/*
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    	-keyout certs/dev.local.io.key -out certs/dev.local.io.crt \
    	-subj "/C=NL/ST=Limburg/L=Geleen/O=Marco Franssen/OU=Development/CN=dev.local.io/emailAddress=marco.franssen@gmail.com"

	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    	-keyout certs/localhost.key -out certs/localhost.crt \
    	-subj "/C=NL/ST=Limburg/L=Geleen/O=Marco Franssen/OU=Development/CN=localhost/emailAddress=marco.franssen@gmail.com"
 */

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

func getCertificate(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			fmt.Printf("%s\nFalling back to Letsencrypt\n", err)
			return certManager.GetCertificate(hello)
		}
		fmt.Println("Loaded selfsigned certificate.")
		return &certificate, err
	}
}

func main() {

	// Pull config
	runningConfig = GetConf("config.yaml")

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

	if runningConfig.docker {
		// Set container mapping type
		docker.ContainerMapType("HOST","godev.sacredheart.it")

		// Subscribe to docker, convert to route and push to router
		go docker.Subscribe(
			"",
			addContainerRoute,
			removeContainerRoute,
		)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/",router.GenericRequestHandler)

	if runningConfig.secure {

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("webnerd.nl","stuff.tld","dev.webnerd.nl"), // Add allowList here

			Cache:      autocert.DirCache("certs"),
		}

		tlsConfig := certManager.TLSConfig()
		tlsConfig.GetCertificate = getCertificate(&certManager)

		server := http.Server{
			Addr:    ":443",
			Handler: mux,
			TLSConfig: tlsConfig,
		}

		go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
		fmt.Println("Server listening on", server.Addr)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			fmt.Println(err)
		}
	} else { // Non secure version
		http.HandleFunc("/",router.GenericRequestHandler)
		http.ListenAndServe(":80", nil)
	}

	// server.ListenAndServeTLS("","")
	// http.HandleFunc("/",router.GenericRequestHandler)
	// certManager.HTTPHandler(nil)
	// http.ListenAndServe(":80",nil)

	// logger.Info("Server is starting on "+host)
	// http.HandleFunc("/",router.GenericRequestHandler)

	// Wrapped in logger.Fatal in case the listenAndServe call ever fails
	// logger.Fatal(http.ListenAndServe(host,nil))
}