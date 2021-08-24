package router

import (
	"bytes"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var errorPage string
var welcomePage string

var currentRoute Route

var hosts []string
var routes []Route

func Init() {
	errorPage = LoadInternalAsset("./assets/error.page.html")
	welcomePage = LoadInternalAsset("./assets/index.page.html")

	// Maybe set logger here to?
}

func AddRoute(route Route) {
	routes = append(routes, route)
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

	println("Something something, Error handler")

	var page = make([]byte,len(errorPage))
	copy(page,errorPage)

	currentRouteResp := fmt.Sprintf("%#v",currentRoute)

	page = bytes.ReplaceAll(page, []byte("{{ERROR_MESSAGE}}"), []byte(err.Error()))
	page = bytes.ReplaceAll(page,[]byte("{{MYOPIA_VERSION}}"),[]byte("Myopia V0.1"))
	page = bytes.ReplaceAll(page,[]byte("{{ROUTE_INFO}}"), []byte(currentRouteResp))

	resp.Write(page)
}

func WelcomePage(resp http.ResponseWriter,r *http.Request) {

	var page = make([]byte,len(welcomePage))
	copy(page,welcomePage)

	page = bytes.ReplaceAll(page,[]byte("{{MYOPIA_VERSION}}"),[]byte("Myopia V0.1"))

	resp.Write(page)
}

func GenericRequestHandler(res http.ResponseWriter,req *http.Request) {

	currentRoute = GetRoute(req)

	switch {
		case currentRoute.MapType == MapHost:
			RouteHost(req,res,currentRoute)
		case currentRoute.MapType == MapAsset:
			ServeStatic(req,res,currentRoute)
		case currentRoute.MapType == MapPath:
			RoutePath(req,res,currentRoute)
	}
}

func LoadInternalAsset(file string) string {
	dat,err := ioutil.ReadFile(file)
	if err != nil {
		// logger.Panicln(err)
	}

	return string(dat)
}

func GetRoute(req *http.Request) Route {

	urlPath := strings.TrimSuffix(req.URL.Path,"/")
	host := strings.SplitN(req.Host,":",2)[0]
	port := strings.SplitN(req.Host,":",2)[1]

	fmt.Printf("Path: [%s] Host: [%s]\n",urlPath,host+":"+port)

	// Check if allowed host
	var allowedHost = false
	for _,hostEntry := range hosts {
		if host == hostEntry {
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

func GetRouteType(req *http.Request) (RouteType, Route) {

	// urlPath := strings.TrimSuffix(req.URL.Path,"/")
	// host := strings.SplitN(req.Host,":",2)[0]
	// port := strings.SplitN(req.Host,":",2)[1]

	/*if runningConfig.debug {
		logger.Debug("Started routing for http(s)://"+host + ":" + port + urlPath+"?"+req.URL.RawQuery)
	}*/

	// Check if host or path
	/*for _,item := range runningConfig.hosts {

		// Host To Host mapping must have a protocol :// defined
		if item.source == "http://"+host {
			return MapHost,item
		}

		if strings.HasPrefix(urlPath,item.source) {
			return MapPath,item
		} else {
			println(item.source," does not start with ",urlPath)
		}
	}

	// Check local asset
	for _,item := range runningConfig.assets {
		if strings.Contains(item.url,urlPath) {
			return MapAsset,defaultRoute
		}
	}*/

	// logger.Warning("Cant map a route for "+urlPath)
	return MapHost,defaultRoute
}

func ServeStatic(req *http.Request,res http.ResponseWriter,route Route) {
	//
}
