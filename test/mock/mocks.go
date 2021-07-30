package mock

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"net/http"
)

type MockClient struct{}

var (
	ValidContent   = []byte("this is the response from a successful request")
	ValidContent2  = []byte("this is the response from another successful request")
	ValidBody      = ioutil.NopCloser(bytes.NewReader(ValidContent))
	ValidBody2     = ioutil.NopCloser(bytes.NewReader(ValidContent2))
	ValidChecksum  = md5.Sum(ValidContent)
	ValidChecksum2 = md5.Sum(ValidContent2)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	switch req.URL.String() {
	case "http://success.com":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ValidBody,
		}, nil
	case "http://success2.com":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ValidBody2,
		}, nil
	case "http://statuscode.com":
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ValidBody,
		}, nil
	case "http://":
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ValidBody,
		}, nil
	default:
		return nil, errors.New("fatal error")
	}
}
