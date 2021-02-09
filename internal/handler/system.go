package handler

import (
	"go.uber.org/zap"
	"net/http"
)

//System handler object
type System struct {
	log *zap.Logger
}

//NewSystemHandler return new dependencies handler instance
func NewSystemHandler(zlog *zap.Logger) *System {
	return &System{log: zlog}
}

//Ping health handler
func (handler System) Ping(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		handler.log.Error("failed to write ping response")
	}
}
