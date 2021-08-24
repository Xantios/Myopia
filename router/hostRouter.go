package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RouteHost(req *http.Request,res http.ResponseWriter,route Route) {

	destinationUrl,_ := url.Parse(route.Destination)
	reverseProxy := httputil.NewSingleHostReverseProxy(destinationUrl)

	// The request we are going to send to the server
	reverseProxy.Director = func(req *http.Request) {

		// main.logger.Debug("Director")

		// Be a good person, make sure X-Proxy headers are set correctly.
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", destinationUrl.Host)
		req.URL.Scheme = "http"
		req.URL.Host = destinationUrl.Host
	}

	reverseProxy.ModifyResponse = func(resp *http.Response) error {

		// if debug {

			//logger.Debug("Response from host")
			//logger.Debug("\t"+resp.Proto+" "+resp.Status)
			//logger.Debug("")

			//for key,_ := range resp.Header {
				// logger.Debug("\t"+key,":",resp.Header.Get(key))
			// }
		//}

		// No errors ! Let's go!
		return error(nil)
	}

	reverseProxy.ErrorHandler = GenericErrorHandler

	reverseProxy.ServeHTTP(res,req)
}
