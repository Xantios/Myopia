package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var currentRoute ConfigItem

func GenericErrorHandler(resp http.ResponseWriter,r *http.Request,err error) {
	logger.Debug("Err!")
	logger.Error("Some server kicked the bucket here! ",err)

	var page = make([]byte,len(errorPage))
	copy(page,errorPage)

	currentRouteResp := fmt.Sprintf("%#v",currentRoute)

	page = bytes.ReplaceAll(page, []byte("{{ERROR_MESSAGE}}"), []byte(err.Error()))
	page = bytes.ReplaceAll(page,[]byte("{{MYOPIA_VERSION}}"),[]byte("Myopia V0.1"))
	page = bytes.ReplaceAll(page,[]byte("{{ROUTE_INFO}}"), []byte(currentRouteResp))

	resp.Write(page)
}

func GenericRequestHandler(res http.ResponseWriter,req *http.Request) {

	routeType, route := GetRouteType(req)
	currentRoute = route

	switch {
		case routeType == MapHost:
			RouteHost(req,res,route)
		case routeType == MapAsset:
			ServeStatic(req,res,route)
		case routeType == MapPath:
			RoutePath(req,res,route)
	}
}

func LoadInternalAsset(file string) string {
	dat,err := ioutil.ReadFile(file)
	if err != nil {
		logger.Panicln(err)
	}

	return string(dat)
}

func GetRouteType(req *http.Request) (RouteType, ConfigItem) {

	urlPath := req.URL.Path
	host := strings.SplitN(req.Host,":",2)[0]
	port := strings.SplitN(req.Host,":",2)[1]

	if runningConfig.debug {
		logger.Debug("Started routing for http(s)://"+host + ":" + port + urlPath+"?"+req.URL.RawQuery)
	}

	// Check if host or path
	for _,item := range runningConfig.hosts {

		// Host To Host mapping must have a protocol :// defined
		if item.source == "http://"+host {
			return MapHost,item
		}

		if strings.Contains(item.source,urlPath) {
			return MapPath,item
		}
	}

	// Check local asset
	for _,item := range runningConfig.assets {
		if strings.Contains(item.url,urlPath) {
			return MapAsset,defaultRoute
		}
	}

	return MapHost,defaultRoute
}

func RoutePath(req *http.Request,res http.ResponseWriter,route ConfigItem) {
	println("Path stuff")
	//
}

func ServeStatic(req *http.Request,res http.ResponseWriter,route ConfigItem) {
	//
}