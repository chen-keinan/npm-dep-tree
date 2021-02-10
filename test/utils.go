package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func requestFunc(handlerFunc func(r http.ResponseWriter, res *http.Request), req *http.Request, target interface{}) (int, error) {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	var err error
	handler.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK {
		if target != nil {
			err = ResJSONToObject(rr, target)
			if err != nil {
				return 0, err
			}
		}
	}
	// Check the status code is what we expect.
	return rr.Code, err
}

//PostWithContext post method for testing
func PostWithContext(ctx context.Context, handlerFunc func(r http.ResponseWriter, res *http.Request), url string, target interface{}, reader io.Reader) (int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, reader)
	if err != nil {
		return 0, err
	}
	return requestFunc(handlerFunc, req, target)

}

//ResJSONToObject map json to Object
func ResJSONToObject(res *httptest.ResponseRecorder, target interface{}) error {
	err := json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		return err
	}
	return nil
}

//InvokeRequest http requests
func InvokeRequest(request *http.Request, handle func(w http.ResponseWriter, r *http.Request), path string, target interface{}) (int, error) {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r)
	})
	m := mux.NewRouter()
	m.HandleFunc(path, f).Methods(request.Method)

	response := httptest.NewRecorder()
	m.ServeHTTP(response, request)
	if response == nil {
		return 0, fmt.Errorf("Failed to get response")
	}
	switch target {
	case target.(string):
		target = response.Body.String()
	default:
		if target != nil {
			err := ResJSONToObject(response, &target)
			if err != nil {
				return 0, err
			}
		}
	}
	return response.Code, nil
}

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

//AddQueryParam add query param to request
func AddQueryParam(request *http.Request, params map[string]string) {
	q := request.URL.Query()
	for pKey, pValue := range params {
		q.Add(pKey, pValue)
	}
	request.URL.RawQuery = q.Encode()
}

//AddHeader add header to request
func AddHeader(request *http.Request, params map[string]string) {
	h := request.Header
	for pKey, pValue := range params {
		h.Add(pKey, pValue)
	}
}

//ReadAppdata read test data
func ReadAppdata(path string) string {
	b, err := ioutil.ReadFile((path))
	if err != nil {
		return ""
	}
	return string(b)
}
