/*
Run this API integration test using the test utility option.
 */
package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"os"
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

var PROTOCOL="http"
var HOST="localhost"
var PORT="9090"

func main() {
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/"},
		&CorrectResponse{http.StatusOK, ".*plugin.*"} )
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/api/public"},
		&CorrectResponse{http.StatusOK, ".*public.*"} )
	runTest(&TestRequest{"GET", PROTOCOL, HOST, PORT, "/api/auth"},
		&CorrectResponse{http.StatusUnauthorized, ".*"} )

	fmt.Println("tests passed...")
}

// Test utility functions
func runTest(request *TestRequest, correct *CorrectResponse) {
	// Prepare request
	url := fmt.Sprintf("%s://%s:%s%s", request.Protocol, request.Host, request.Port, request.Route)
	req, err := http.NewRequest(request.Method, url, nil)
	if err != nil { failTest("build request failed...") }

	// Run request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { failTest("send request failed...") }
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { failTest("parsing response failed...") }

	// Test the response
	if resp.StatusCode != correct.ResponseCode { failTest("incorrect response code...") }
	matched, err := regexp.MatchString(correct.ContentRegexp, string(content))
	if !matched || err != nil { failTest("contents check failed...") }

	fmt.Println("passed: ", request)
}

func failTest(failMessage string) {
	fmt.Println(failMessage)
	os.Exit(-1)
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
