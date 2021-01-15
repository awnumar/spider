package main

import (
	"fmt"
	"io"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/net/html"
)

type pageLinks struct {
	pageAddr *url.URL
	links    []*url.URL
	errors   []*parsingError
}

type parsingError struct {
	link string
	err  error
}

func crawl(rootURL string) error {
	rootAddr, err := url.Parse(rootURL)
	if err != nil {
		return fmt.Errorf("error parsing address %s: %w", rootURL, err)
	}

	seen := map[string]bool{}

	queue := make(chan *url.URL, 256)
	defer close(queue)

	queue <- rootAddr

	for address := range queue {
		if !seen[address.String()] && address.Host == rootAddr.Host {
			seen[address.String()] = true

			if err := parsePage(address, queue); err != nil {
				return err
			}
		}
		if len(queue) == 0 {
			break
		}
	}

	fmt.Print("\nDiscovered endpoints:\n\n")
	for k := range seen {
		fmt.Println(k)
	}

	return nil
}

func startRoutine(address *url.URL) {

}

func parsePage(address *url.URL) (*pageLinks, error) {
	resp, err := retryablehttp.Get(address.String())
	if err != nil {
		return nil, fmt.Errorf("error performing GET request to %s: %w", address, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error: status %d from %s: %s", resp.StatusCode, address, resp.Status)
	}

	results := &pageLinks{
		pageAddr: address,
		links:    []*url.URL{},
		errors:   []*parsingError{},
	}

	page := html.NewTokenizer(resp.Body)

	for {
		switch page.Next() {

		case html.StartTagToken:
			token := page.Token()
			if token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						parsedLink, err := url.Parse(attribute.Val)
						if err != nil {
							results.errors = append(results.errors, &parsingError{
								link: attribute.Val,
								err:  fmt.Errorf("skip: could not parse address %s: %w", attribute.Val, err),
							})
							break
						}
						if parsedLink.Host == "" { // convert relative links to absolute
							parsedLink.Scheme = address.Scheme
							parsedLink.User = address.User
							parsedLink.Host = address.Host
						}
						if parsedLink.Host == address.Host {
							results.links = append(results.links, parsedLink)
						}
						break
					}
				}
			}

		case html.ErrorToken:
			if page.Err() != io.EOF {
				return results, fmt.Errorf("error tokenizing html: %w", page.Err())
			}
			return results, nil
		}
	}
}
