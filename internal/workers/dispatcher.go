package workers

import (
	"github.com/chen-keinan/npm-dep-tree/internal/router/middleware"
	"github.com/chen-keinan/npm-dep-tree/internal/workers/processing"
)

//WorkerDispatcher object
type WorkerDispatcher struct {
}

//NewWorkerDispatcher return new worker dispatcher
func NewWorkerDispatcher() *WorkerDispatcher {
	return &WorkerDispatcher{}
}

//Work worker to invoke processing function and call designated handler
func (wd WorkerDispatcher) Work(jobFunc func(msg middleware.RequestProcessor), prChan chan middleware.RequestProcessor) {
	for msg := range prChan {
		jobFunc(msg)
	}
}

//InvokeProcessingWorkers initiate 50 worker to process requests in async
func (wd WorkerDispatcher) InvokeProcessingWorkers(pr chan middleware.RequestProcessor) {
	for index := 0; index < 50; index++ {
		worker := processing.NewReqProcessingWorker()
		// invoke request processing workers
		go wd.Work(worker.ProcessRequests, pr)
	}
}
