package router

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RouteHost(req *http.Request,res http.ResponseWriter,route Route) {

	destinationUrl,_ := url.Parse(route.Destination)
	reverseProxy := httputil.NewSingleHostReverseProxy(destinationUrl)

	// Ignore downstream SSL cert errors
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// The request we are going to send to the server
	reverseProxy.Director = func(req *http.Request) {

		// destinationUrl.Host
		compoundPath := destinationUrl.Path+req.URL.Path

		// Be a good person, make sure X-Proxy headers are set correctly.
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", req.Host)
		req.URL.Scheme = destinationUrl.Scheme
		req.URL.Host = destinationUrl.Host
		req.Host = destinationUrl.Host
		req.URL.Path = compoundPath

		if debug {
			LogRequest(req)
		}
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {

		if debug {
			LogResponse(resp)
		}

		// No errors ! Let's go!
		return error(nil)
	}

	reverseProxy.ErrorHandler = GenericErrorHandler
	reverseProxy.ServeHTTP(res,req)
}
