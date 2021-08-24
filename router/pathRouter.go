package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RoutePath(req *http.Request,res http.ResponseWriter,route ConfigItem) {

	destinationUrl,_ := url.Parse(route.destination)
	sourceUrl,_ := url.Parse("http://"+req.Host+route.source)

	reverseProxy := httputil.NewSingleHostReverseProxy(destinationUrl)

	// The request we are going to send to the server
	reverseProxy.Director = func(req *http.Request) {

		logger.Debug("Director")

		// Be a good person, make sure X-Proxy headers are set correctly.
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", sourceUrl.Host)
		req.URL.Scheme = "http"
		// req.URL.Host = destinationUrl.Host
		req.URL.Host = "KOFFIEKAAS.TLD"
		req.URL.Path = "/api/characters"

		dump,err := httputil.DumpRequestOut(req,true)
		if err != nil {
			logger.Error("Request broken!",err)
		}

		fmt.Printf("%q",dump)


	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {

		if runningConfig.debug {

			logger.Debug("Response from host")
			logger.Debug("\t"+resp.Proto+" "+resp.Status)
			logger.Debug("")

			for key,_ := range resp.Header {
				logger.Debug("\t"+key,":",resp.Header.Get(key))
			}
		}

		// No errors ! Let's go!
		return error(nil)
	}

	reverseProxy.ErrorHandler = GenericErrorHandler
	reverseProxy.ServeHTTP(res,req)
}