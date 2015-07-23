package main

import (
	"fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
//func Crawl(url string, depth int, fetcher Fetcher) {
//	// TODO: Fetch URLs in parallel.
//	// TODO: Don't fetch the same URL twice.
//	// This implementation doesn't do either:
//	if depth <= 0 {
//		return
//	}
//	body, urls, err := fetcher.Fetch(url)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("found: %s %q\n", url, body)
//	for _, u := range urls {
//		Crawl(u, depth-1, fetcher)
//	}
//	return
//}


func Krawl(url string, fetcher Fetcher, Urls chan []string) {
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("found: %s %q\n", url, body)
	}
	Urls <- urls
}

func Crawl(url string, depth int, fetcher Fetcher) {
	store := make(map[string]bool)
	Urls := make(chan []string)
	go Krawl(url, fetcher, Urls)
	band := 1
	store[url] = true // init for level 0 done
	for i := 0; i < depth; i++ {
		for band > 0 {
			band--
			next := <- Urls
			for _, url := range next {
				if _, done := store[url] ; !done {
					store[url] = true
					band++
					go Krawl(url, fetcher, Urls)
				}
			}
		}
	}
	return
}


func main() {
	Crawl("http://golang.org/", 4, fetcher)
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
