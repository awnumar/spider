package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/asaskevich/govalidator"
)

func main() {
	var rootURL string
	var frequency float64

	flag.StringVar(&rootURL, "url", "", "spider will start crawling at this URL")
	flag.Float64Var(&frequency, "frequency", 5, "specify the maximum number of requests per second")
	flag.Parse()

	if rootURL == "" {
		fmt.Println("error: must provide url parameter")
		flag.PrintDefaults()
		return
	}
	if !govalidator.IsURL(rootURL) {
		fmt.Println("error: url parameter must be a valid URL")
		return
	}

	parsedURL, err := url.Parse(rootURL)
	if err != nil {
		fmt.Printf("error parsing url: %s\n", err)
		return
	}

	delayBetweenRequests := time.Duration((1.0/frequency)*float64(time.Second)) * time.Nanosecond

	if err := crawl(parsedURL, delayBetweenRequests); err != nil {
		fmt.Println(err)
	}
}
