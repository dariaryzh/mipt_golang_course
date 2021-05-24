package httpfetch

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

type Request struct {
	Method string
	URL    string
	Body   []byte
}

type Result struct {
	StatusCode int
	Error      error
}

func FetchAll(c *http.Client, requests []Request) []Result {
	var wg sync.WaitGroup
	wg.Add(len(requests))
	resCh := make(chan Result)

	for i, v := range requests {
		log.Printf("Got a request %d", i)
		go makerequest(c, &wg, v, resCh)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	results := make([]Result, 0, len(requests))
	for r := range resCh {
		results = append(results, r)
	}
	return results
}

func makerequest(c *http.Client, wg *sync.WaitGroup, r Request, resCh chan Result) {
	defer wg.Done()

	req, err := http.NewRequest(r.Method, r.URL, bytes.NewReader(r.Body))

	if err != nil {
		return
	}

	var res Result
	resp, err := c.Do(req)
	if err != nil {
		res.Error = err
		resCh <- res
		return
	}

	defer resp.Body.Close()
	res.StatusCode = resp.StatusCode
	resCh <- res
}