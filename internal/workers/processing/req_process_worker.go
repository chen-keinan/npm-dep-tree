package processing

import (
	"github.com/chen-keinan/npm-dep-tree/internal/router/middleware"
)

//ReqProcessingWorker object
type ReqProcessingWorker struct {
	Type string
}

//NewReqProcessingWorker worker to process dependencies tree request
func NewReqProcessingWorker() *ReqProcessingWorker {
	return &ReqProcessingWorker{Type: "RequestProcessing"}
}

//ProcessRequests processing request and invoke dependencies handler
func (rpw ReqProcessingWorker) ProcessRequests(msg middleware.RequestProcessor) {
	msg.Handler.ServeHTTP(msg.ResponseWriter, msg.Request)
	// notify back to middleware that processing has completed
	msg.Completed <- true
}
