package myhttp

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/FdeFabricio/http-md5/test/mock"
)

func init() {
	Client = &mock.MockClient{}
}

var ttValidateURL = []struct {
	name    string
	input   string
	wantErr bool
	wantURL string
}{
	{
		name:    "invalid url",
		input:   " ",
		wantErr: true,
	},
	{
		name:    "valid url without http",
		input:   "google.com",
		wantErr: false,
		wantURL: "http://google.com",
	},
	{
		name:    "valid url with http",
		input:   "http://google.com",
		wantErr: false,
		wantURL: "http://google.com",
	},
}

func TestValidateURL(t *testing.T) {
	for _, tt := range ttValidateURL {
		t.Run(tt.name, func(t *testing.T) {
			newURL, err := validateURL(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("validateURL -> expected error %v, got %v", tt.wantErr, err)
				}
			} else if newURL != tt.wantURL {
				t.Errorf("validateURL -> expected %v, got %v", tt.wantURL, newURL)
			}
		})
	}
}

var ttGetMD5 = []struct {
	name     string
	url      string
	checksum [16]byte
	err      error
}{
	{
		name:     "success returns valid checksum",
		url:      "http://success.com",
		checksum: mock.ValidChecksum,
		err:      nil,
	},
	{
		name:     "request fails",
		url:      "http://fails.com",
		checksum: [16]byte{},
		err:      errors.New("failed to get URL"),
	},
	{
		name:     "request returns status code different than 200",
		url:      "http://statuscode.com",
		checksum: [16]byte{},
		err:      errors.New("failed with StatusCode="),
	},
}

func TestGetMD5(t *testing.T) {
	for _, tt := range ttGetMD5 {
		t.Run(tt.name, func(t *testing.T) {
			checksum, err := getMD5(tt.url)
			if err != nil && tt.err == nil {
				t.Errorf("getMD5 -> returned unexpected error: %s", err)
			} else if err == nil && tt.err != nil {
				t.Errorf("getMD5 -> expected error not returned: %s", tt.err)
			} else if err != nil && tt.err != nil && !strings.Contains(err.Error(), tt.err.Error()) {
				t.Errorf("getMD5 -> expected error like '%s', got '%s'", tt.err, err)
			} else if checksum != tt.checksum {
				t.Errorf("getMD5 -> expected checksum '%x', got '%x'", tt.checksum, checksum)
			}
		})
	}
}

var ttExecute = []struct {
	name   string
	urls   []string
	output string
}{
	{
		name:   "multiple successful requests",
		urls:   []string{"http://success.com", "success2.com"},
		output: fmt.Sprintf("http://success.com %x\nhttp://success2.com %x", mock.ValidChecksum, mock.ValidChecksum2),
	},
	{
		name:   "invalid url request",
		urls:   []string{""},
		output: "[ERROR] invalid URL 'http://'",
	},
	{
		name:   "request with error",
		urls:   []string{"http://statuscode.com"},
		output: "[ERROR] request http://statuscode.com failed with StatusCode=400",
	},
}

// test what is printed (sent to stdout)
func TestExecute(t *testing.T) {
	for _, tt := range ttExecute {
		// after a body is read using `io.ReadAll` it needs to be recreated
		mock.ValidBody = ioutil.NopCloser(bytes.NewReader(mock.ValidContent))
		mock.ValidBody2 = ioutil.NopCloser(bytes.NewReader(mock.ValidContent2))

		t.Run(tt.name, func(t *testing.T) {
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			Execute(1, tt.urls)

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			if strings.TrimSpace(string(out)) != tt.output {
				t.Errorf("Execute -> expected output \n'%s' got \n'%s'", tt.output, out)
			}
		})
	}
}
