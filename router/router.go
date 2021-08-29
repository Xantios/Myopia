package router

import (
	"bytes"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var errorPage string
var welcomePage string

var currentRoute Route

var hosts []string
var routes []Route

// RouterLog routign specific logging
var RouterLog *logrus.Entry
var debug bool

func Init() {

	// Setup logging
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	RouterLog = logrus.WithFields(logrus.Fields{
		"service": "Router",
	})

	errorPage = LoadInternalAsset("./assets/error.page.html")
	welcomePage = LoadInternalAsset("./assets/index.page.html")
}

func SetDebug(debugValue bool) {
	debug = debugValue
}

func AddRoute(route Route) {
	routes = append(routes, route)
}

func RemoveRoute(name string) {

	// Convert name to index
	matchIndex := -1
	for index,route := range routes {
		if route.Name == name {
			matchIndex = index
		}
	}

	if matchIndex == -1 {
		println("No such route "+name)
		return
	}

	// Move to bottom, slice of last item
	routes[matchIndex] = routes[len(routes)-1]
	routes = routes[:len(routes)-1]
}

func AddHost(host string) {
	hosts = append(hosts, host)
}

func PrintRouteTable() {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredDark)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name","Source","Destination"})
	for _,route := range routes {
		// fmt.Printf("%s\t\t%s\n", route.Source,route.Destination)
		t.AppendRow([]interface{}{route.Name,route.Source,route.Destination})
	}
	t.AppendFooter(table.Row{"Routes", len(routes)})
	t.Render()

}

func GenericErrorHandler(resp http.ResponseWriter,r *http.Request,err error) {

	var page = make([]byte,len(errorPage))
	copy(page,errorPage)

	currentRouteResp := fmt.Sprintf("%#v",currentRoute)

	page = bytes.ReplaceAll(page, []byte("{{ERROR_MESSAGE}}"), []byte(err.Error()))
	page = bytes.ReplaceAll(page,[]byte("{{MYOPIA_VERSION}}"),[]byte("Myopia V0.1"))
	page = bytes.ReplaceAll(page,[]byte("{{ROUTE_INFO}}"), []byte(currentRouteResp))

	resp.Write(page)
}

func GenericRequestHandler(res http.ResponseWriter,req *http.Request) {

	currentRoute = GetRoute(req)

	switch {
		case currentRoute.MapType == MapHost:

			if debug {
				RouterLog.Info("Using Host router")
			}

			RouteHost(req,res,currentRoute)

		case currentRoute.MapType == MapPath:

			if debug {
				RouterLog.Info("Using Path router")
			}

			RoutePath(req,res,currentRoute)
	}
}

func LoadInternalAsset(file string) string {
	dat,err := os.ReadFile(file)
	if err != nil {
		RouterLog.Warning("Cant load internal asset ["+file+"]")
		return ""
	}

	return string(dat)
}

func GetRoute(req *http.Request) Route {

	urlPath := strings.TrimSuffix(req.URL.Path,"/")

	// Make sure port is always defined
	if !strings.Contains(req.Host,":") {
		req.Host = req.Host+":80"
	}

	host := strings.SplitN(req.Host,":",2)[0]
	port := strings.SplitN(req.Host,":",2)[1]

	fmt.Printf("Path: [%s] Host: [%s]\n",urlPath,host+":"+port)

	// Check if allowed host
	var allowedHost = false
	for _,hostEntry := range hosts {
		if  strings.HasSuffix(host,hostEntry) {
			allowedHost = true
		}
	}

	if !allowedHost {
		println("Host "+host+" is not allowed")
		return defaultRoute
	}

	// Check if we can map host to host
	var challengeHost = "http://"+host
	for _,route := range routes {
		if route.Source == challengeHost {
			return route
		}
	}

	// Check if we can map path to host
	for _,route := range routes {
		if strings.HasPrefix(urlPath,route.Source) {
			return route
		}
	}

	return defaultRoute
}
