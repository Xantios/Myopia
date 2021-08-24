package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func RoutePath(req *http.Request,res http.ResponseWriter,route Route) {

	destinationUrl,_ := url.Parse(route.Destination)
	sourceUrl,_ := url.Parse("http://"+req.Host+route.Source)

	reverseProxy := httputil.NewSingleHostReverseProxy(destinationUrl)

	// The request we are going to send to the server
	reverseProxy.Director = func(req *http.Request) {

		target,parseError := url.Parse(currentRoute.Destination)

		if parseError != nil {
			println("Cant parse target URL:",currentRoute.Destination)
			return
		}

		// Craft the target URL together by replacing the RoutePath in the request URL
		fixedPath := strings.Replace(req.URL.Path,route.Source,target.Path,1)

		// Be a good person, make sure X-Proxy headers are set correctly.
		req.Header.Add("X-Forwarded-Host",sourceUrl.Host)
		req.Header.Add("X-Origin-Host",  req.Host)
		req.URL.Scheme = "http"
		req.Host = target.Host
		// req.URL.Path = "/api/characters"
		req.URL.Path = fixedPath
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		dump,err := httputil.DumpRequestOut(req,true)
		if err != nil {
			println("Request broken",err.Error())
		}
		fmt.Printf("%#q",dump)

		// println(target.Path)
		// fmt.Fprint(res,string(dump))
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {

		println("Response from server:")
		println(resp.Proto+" "+resp.Status)
		for key,_ := range resp.Header {
			println(key+": "+resp.Header.Get(key))
		}

		// No errors ! Let's go!
		return error(nil)
	}

	reverseProxy.ErrorHandler = GenericErrorHandler
	reverseProxy.ServeHTTP(res,req)
}