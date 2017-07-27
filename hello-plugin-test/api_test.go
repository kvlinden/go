/*
Run this API integration test using the test utility option.
 */
package main

import (
	"net/http"
	"testing"
	"fmt"
	"io/ioutil"
	"regexp"
)

type TestRequest struct {
	Method      string
	Protocol    string
	Host        string
	Port        string
	Route       string
}
type CorrectResponse struct {
	ResponseCode int
	ContentRegexp string
}

// The tests
func TestHome(test *testing.T) {
	testRequest := &TestRequest{"GET", "http", "localhost", "8080", "/"}
	correctResponse := &CorrectResponse{http.StatusOK, ".*plugin.*"}
	runTest(testRequest, correctResponse, test)
}

func TestPublic(test *testing.T) {
	testRequest := &TestRequest{"GET", "http", "localhost", "8080", "/api/public"}
	correctResponse := &CorrectResponse{http.StatusOK, ".*public.*"}
	runTest(testRequest, correctResponse, test)
}

func TestAdmin(test *testing.T) {
	testRequest := &TestRequest{"GET", "http", "localhost", "8080", "/api/public"}
	correctResponse := &CorrectResponse{http.StatusOK, ".*public.*"}
	runTest(testRequest, correctResponse, test)
}

// Test utility functions
func runTest(request *TestRequest, correct *CorrectResponse, test *testing.T) {
	// Prepare request
	url := fmt.Sprintf("%s://%s:%s%s", request.Protocol, request.Host, request.Port, request.Route)
	req, err := http.NewRequest(request.Method, url, nil)
	if err != nil { failTest(test, "build request failed...") }

	// Run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { failTest(test,"send request failed...") }
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { failTest(test, "parsing response failed...") }

	// Test the response
	if resp.StatusCode != correct.ResponseCode { failTest(test,"incorrect response code...") }
	matched, err := regexp.MatchString(correct.ContentRegexp, string(content))
	if !matched || err != nil { failTest(test, "contents check failed...") }
}

func failTest(test *testing.T, failMessage string) {
	test.Log(failMessage)
	test.Fail()
}



// https://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest
//func runTest(test *TestRequest) {
//	//apiUrl := "https://api.com"
//	//resource := "/user/"
//	//data := url.Values{}
//	//data.Set("name", "foo")
//	//data.Add("surname", "bar")
//
//	fmt.Println("***", test)
//
//	//u, _ := url.ParseRequestURI()
//	//u.Path = test.Route
//	//urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"
//	//
//	//fmt.Println("***", urlStr, "***")
//
//	client := &http.Client{}
//	req, _ := http.NewRequest("GET", fmt.Sprintf("%s://%s:%s%s", test.Protocol, test.Host, test.Port), nil)
//	//r, _ := http.NewRequest("GET", urlStr, bytes.NewBufferString(data.Encode()))
//	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
//	//r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
//	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
//
//	resp, _ := client.Do(req)
//	fmt.Println(resp.Status)
//}



