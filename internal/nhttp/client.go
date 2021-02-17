package nhttp

import "net/http"

//HClient interface expose one get method
//client.go
//go:generate mockgen -destination=../mocks/mock_HClient.go -package=mocks . HClient
type HClient interface {
	Get(url string) (*http.Response, error)
}

//NpmHTTPClient object for communicating with npm registry
type NpmHTTPClient struct {
	hclient *http.Client
}

//NewNpmHTTPClient instantiate new Npm http client object
func NewNpmHTTPClient() HClient {
	return &NpmHTTPClient{hclient: &http.Client{}}
}

//Get HTTP Method implementation
func (nhc NpmHTTPClient) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return nhc.hclient.Do(request)
}
