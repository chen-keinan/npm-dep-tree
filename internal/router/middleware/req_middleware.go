package middleware

import (
	"net/http"
	"strings"
)

//RequestProcessor request limit msg with handler and request data
type RequestProcessor struct {
	Handler        http.Handler
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Completed      chan bool
}

//NewRateLimitChan return rate limit chan
func NewRateLimitChan() chan RequestProcessor {
	return make(chan RequestProcessor)
}

//RequestLimitMiddleware perform request limit for dependency request and process it in async
func RequestLimitMiddleware(h http.Handler, ProcessingMsg chan RequestProcessor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/api/v1/package-dependencies/") {
			// request processing message
			procMsg := RequestProcessor{Handler: h, ResponseWriter: w, Request: req, Completed: make(chan bool)}
			go func() {
				ProcessingMsg <- procMsg
			}()
			// wait for handler to complete
			<-procMsg.Completed
			return
		}
		// continue processing  request
		h.ServeHTTP(w, req)
	})
}
