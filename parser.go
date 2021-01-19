package main

import (
	"fmt"
	"io"
	"net/url"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/html"
)

type pageLinks struct {
	addr   *url.URL
	links  []*url.URL
	errors []*pageError
}

type pageError struct {
	addr string
	err  error
}

func extractLinks(r io.Reader, result *pageLinks) *pageLinks {
	page := html.NewTokenizer(r)

	for {
		switch page.Next() {

		case html.StartTagToken:
			token := page.Token()
			if token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						parsedLink, err := url.Parse(attribute.Val)
						if err != nil {
							result.errors = append(result.errors, &pageError{
								addr: attribute.Val,
								err:  fmt.Errorf("error: could not parse address %s: %w", attribute.Val, err),
							})
							break // skip this link and continue parsing
						}
						parsedLink.Fragment = ""
						if parsedLink.Host == "" { // convert relative links to absolute
							parsedLink.Scheme = result.addr.Scheme
							parsedLink.User = result.addr.User
							parsedLink.Host = result.addr.Host
						}
						if parsedLink.Host == result.addr.Host && govalidator.IsURL(parsedLink.String()) && !contains(result.links, parsedLink) {
							result.links = append(result.links, parsedLink)
						}
						break
					}
				}
			}

		case html.ErrorToken:
			if page.Err() != io.EOF {
				result.errors = append(result.errors, &pageError{
					addr: result.addr.String(),
					err:  fmt.Errorf("error parsing html: %w", page.Err()),
				})
			}
			return result
		}
	}
}

func contains(list []*url.URL, item *url.URL) bool {
	for _, v := range list {
		if v.String() == item.String() {
			return true
		}
	}
	return false
}
