package main

import (
	"flag"
	"fmt"

	"github.com/asaskevich/govalidator"
)

func main() {
	var url string

	flag.StringVar(&url, "url", "", "spider will start crawling at this URL")
	flag.Parse()

	if url == "" {
		fmt.Println("error: must provide url parameter")
		flag.PrintDefaults()
		return
	}
	if !govalidator.IsURL(url) {
		fmt.Println("error: url parameter must be a valid URL")
		flag.PrintDefaults()
		return
	}

	if err := crawl(url); err != nil {
		fmt.Println(err)
	}
}
