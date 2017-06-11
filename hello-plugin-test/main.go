/*
Run this API integration test using the test utility option.
*/
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type TestRequest struct {
	Method   string
	Protocol string
	Host     string
	Port     string
	Route    string
}
type CorrectResponse struct {
	ResponseCode  int
	ContentRegexp string
}

var PROTOCOL = "http"
var HOST = "localhost"
var PORT = "9090"

func main() {
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/"},
		&CorrectResponse{http.StatusOK, ".*plugin.*"})
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/api/public"},
		&CorrectResponse{http.StatusOK, ".*public.*"})
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/api/auth"},
		&CorrectResponse{http.StatusUnauthorized, ".*"})

	testHttpInfoPassage()

	fmt.Println("tests passed...")
}

// Test utility functions
func runTest(request *TestRequest, correct *CorrectResponse) {
	// Prepare request
	url := fmt.Sprintf("%s://%s:%s%s", request.Protocol, request.Host, request.Port, request.Route)
	req, err := http.NewRequest(request.Method, url, nil)
	if err != nil {
		failTest("build request failed...")
	}

	// Run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		failTest("send request failed...")
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		failTest("parsing response failed...")
	}

	// Test the response
	if resp.StatusCode != correct.ResponseCode {
		failTest("incorrect response code...")
	}
	matched, err := regexp.MatchString(correct.ContentRegexp, string(content))
	if !matched || err != nil {
		failTest("contents check failed...")
	}

	fmt.Println("passed: ", request)
}

func failTest(failMessage string) {
	fmt.Println(failMessage)
	os.Exit(-1)
}

func testHttpInfoPassage() {
	// Prepare request (with post data and a header key value)
	url := "http://localhost:9090/"
	postData := []byte(`{"key":"post data"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("header-key", "header value")
	if err != nil {
		failTest("build request failed...")
	}

	// Run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		failTest("send request failed...")
	}

	// Process response
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		failTest("parsing response failed...")
	}
	fmt.Println("***", "content: ", string(content))
	fmt.Println("***", "code: ", resp.Status)
	fmt.Println("***", "a sample response header value - content-length: ", resp.Header.Get("Content-Length"))

}
