# Spider

An aggressive, multi-threaded spider for discovering all reachable HTTP endpoints on a host.

## Usage

1. `go build .`
2. `./spider -help`

## Notes

- Some sites or WAFs may block this tool for sending too many requests too quickly. A mitigation strategy is to use the `-frequency` flag to limit the rate of requests.
- Side effect free functions are unit tested but integration tests would also be useful if this was to be deployed in production.
- The current status indication features are not great. There's a debug output that shows requests being made and any errors, and results are printed after crawling is finished. A better way would be to refresh the screen with information such as links visited, links in queue, unique links discovered, etc. Also, a sitemap could be generated instead of links being printed out.
