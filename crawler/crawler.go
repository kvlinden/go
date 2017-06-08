package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

const (
	FETCHING = iota
	FETCHED
	FAILED
)

type URLList struct {
	v             map[string]int
	fetchingCount int
	mux           sync.Mutex
}

func (urlList *URLList) Append(newUrl string) string {
	urlList.mux.Lock()
	_, exists := urlList.v[newUrl]
	if !exists {
		urlList.v[newUrl] = FETCHING
		urlList.fetchingCount = urlList.fetchingCount + 1
	} else {
		newUrl = ""
	}
	urlList.mux.Unlock()
	return newUrl
}

func (urlList *URLList) MarkFetched(fetchedUrl string) {
	urlList.mux.Lock()
	urlList.v[fetchedUrl] = FETCHED
	urlList.fetchingCount = urlList.fetchingCount - 1
	urlList.mux.Unlock()
}

func (urlList *URLList) MarkFailed(fetchedUrl string) {
	urlList.mux.Lock()
	urlList.v[fetchedUrl] = FAILED
	urlList.fetchingCount = urlList.fetchingCount - 1
	urlList.mux.Unlock()
}

func (urlList *URLList) CheckStop(ch chan string) {
	urlList.mux.Lock()
	if urlList.fetchingCount == 0 {
		close(ch)
	}
	urlList.mux.Unlock()
}

func Crawl(url string, depth int, fetcher Fetcher, urlList *URLList, ch chan string) {
	defer urlList.CheckStop(ch)
	if depth <= 0 {
		urlList.MarkFailed(url)
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		urlList.MarkFailed(url)
		return
	}
	ch <- fmt.Sprintf("found: %s %q in %v", url, body, urls)
	urlList.MarkFetched(url)
	for _, u := range urls {
		if urlList.Append(u) != "" {
			go Crawl(u, depth-1, fetcher, urlList, ch)
		}
	}
	return
}

func main() {
	urlList := &URLList{v: make(map[string]int)}
	ch := make(chan string, 10)
	Crawl(urlList.Append("http://golang.org/"), 4, fetcher, urlList, ch)
	for i := range ch {
		fmt.Println(i)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
