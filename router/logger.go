package router

import "net/http"

func LogResponse(resp *http.Response) {
	RouterLog.Debug("Response from host")
	RouterLog.Debug("\t"+resp.Proto+" "+resp.Status)
	RouterLog.Debug("")

	for key,_ := range resp.Header {
		RouterLog.Debug("\t"+key,":",resp.Header.Get(key))
	}
}

func LogRequest(req *http.Request) {
	RouterLog.Debug("Request to host ["+req.Host+"]")
	RouterLog.Debug("\t"+req.Method+" "+req.RequestURI+" "+req.Proto)
	RouterLog.Debug("Host: "+req.Host)

	for key,value := range req.Header {
		RouterLog.Debug("\t"+key+": "+value[0])
	}
}
