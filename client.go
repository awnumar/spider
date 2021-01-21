package main

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type client struct {
	c     *retryablehttp.Client
	reqs  chan *request
	delay time.Duration
}

type request struct {
	addr string
	resp chan *response
}

type response struct {
	resp *http.Response
	err  error
}

func newClient(delay time.Duration) *client {
	cl := &client{
		c:     retryablehttp.NewClient(),
		reqs:  make(chan *request),
		delay: delay,
	}

	go func() {
		last := time.Now()
		for cmd := range cl.reqs {
			if time.Since(last) < cl.delay {
				time.Sleep(cl.delay - time.Since(last))
			}
			last = time.Now()

			req, err := retryablehttp.NewRequest(http.MethodGet, cmd.addr, nil)
			if err != nil {
				cmd.resp <- &response{err: err}
			} else {
				resp, err := retryablehttp.NewClient().Do(req)
				cmd.resp <- &response{resp, err}
			}

			close(cmd.resp)
		}
	}()

	return cl
}

func (c *client) Get(addr string) <-chan *response {
	req := &request{
		addr: addr,
		resp: make(chan *response),
	}
	c.reqs <- req
	return req.resp
}
