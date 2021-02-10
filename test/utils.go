package test

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)

//InvokeRequestWithResponse http requests
func InvokeRequestWithResponse(request *http.Request, handle func(w http.ResponseWriter, r *http.Request), path string) (*httptest.ResponseRecorder, error) {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r)
	})
	m := mux.NewRouter()
	response := httptest.NewRecorder()
	m.HandleFunc(path, f).Methods(request.Method)
	m.ServeHTTP(response, request)
	return response, nil
}
