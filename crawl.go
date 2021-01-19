package main

import (
	"fmt"
	"net/url"
	"time"
)

func crawl(rootAddr *url.URL, delayBetweenRequests time.Duration) error {
	fmt.Printf("Crawling %s with a minimum delay of %s between requests.\n\n", rootAddr, delayBetweenRequests.String())
	client := newClient(delayBetweenRequests)

	done := map[string]*pageLinks{}

	// results from crawling routines received here
	results := make(chan *pageLinks)

	// change in the number of active routines received here
	delta := make(chan int)

	// queue of yet to crawl pages wait here
	queue := make(chan *url.URL)

	go func() {
		numberOfRoutines := 0
		for d := range delta {
			numberOfRoutines += d
			if numberOfRoutines == 0 {
				close(queue)
			}
		}
	}()

	go func() {
		for endpoint := range queue {
			go func(addr *url.URL) {
				delta <- 1
				results <- parsePage(addr, client)

			}(endpoint)
		}
		close(delta)
		close(results)
	}()

	queue <- rootAddr

	for r := range results {
		done[r.addr.String()] = r

		for _, link := range r.links {
			if done[link.String()] == nil {
				done[link.String()] = &pageLinks{addr: link}
				queue <- link
			}
		}

		delta <- -1
	}

	for crawled := range done {
		if len(done[crawled].links) != 0 {
			fmt.Printf("\n\n:: Crawled page at URL %s, found these links:\n\n", crawled)
			for _, result := range done[crawled].links {
				fmt.Printf("+ %s\n", result.String())
			}
			if len(done[crawled].errors) != 0 {
				fmt.Print("\n with these errors encountered:\n\n")
				for _, err := range done[crawled].errors {
					fmt.Printf("! link=%s, err=%s\n", err.addr, err.err)
				}
			}
		} else {
			if len(done[crawled].errors) != 0 {
				fmt.Printf("\n\nErrors encountered while crawling %s:\n", done[crawled].addr.String())
				for _, err := range done[crawled].errors {
					fmt.Printf("! link=%s, err=%s\n", err.addr, err.err)
				}
			}
		}
	}

	fmt.Printf("\n\n:: Found these unique links on %s:\n\n", rootAddr)
	for crawled := range done {
		fmt.Println("+", crawled)
	}

	return nil
}

func parsePage(address *url.URL, c *client) *pageLinks {
	result := &pageLinks{
		addr:   address,
		links:  []*url.URL{},
		errors: []*pageError{},
	}

	maybeResponse := <-c.Get(address.String())
	if maybeResponse.err != nil {
		result.errors = append(result.errors, &pageError{
			addr: address.String(),
			err:  maybeResponse.err,
		})
		return result
	}
	resp := maybeResponse.resp
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result.errors = append(result.errors, &pageError{
			addr: address.String(),
			err:  fmt.Errorf("error: status %d from %s: %s", resp.StatusCode, address, resp.Status),
		})
		return result
	}

	return extractLinks(resp.Body, result)
}
