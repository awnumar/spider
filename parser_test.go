package main

import (
	"bytes"
	"net/url"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	addr, err := url.Parse("https://spacetime.dev")
	if err != nil {
		t.Error(err)
	}

	links := extractLinks(bytes.NewBuffer([]byte(sampleHTML)), &pageLinks{addr: addr})

	if addr.String() != links.addr.String() {
		t.Error("function mutated address argument")
	}

	if len(links.errors) != 0 {
		t.Error("didn't expect any errors, got", links.errors)
	}

	if !equal(links.links, containedLinks) {
		t.Error("extracted links different than expected; got", links.links, "ref", containedLinks)
	}
}

func TestContains(t *testing.T) {
	addr, err := url.Parse("https://spacetime.dev")
	if err != nil {
		t.Error(err)
	}

	list := []*url.URL{addr}

	if !contains(list, addr) {
		t.Error("contained returned false when it should have returned true")
	}

	list = []*url.URL{}

	if contains(list, addr) {
		t.Error("contained returned true when it should have returned false")
	}
}

var containedLinks = []string{
	"https://spacetime.dev/",
	"https://spacetime.dev/posts",
	"https://spacetime.dev/feed.xml",
	"https://spacetime.dev/cdn-cgi/l/email-protection",
	"https://spacetime.dev/public-key.pgp",
	"https://spacetime.dev/rosen-censorship-resistant-proxy-tunnel",
	"https://spacetime.dev/plausibly-deniable-encryption",
	"https://spacetime.dev/memory-retention-attacks",
	"https://spacetime.dev/encrypting-secrets-in-memory",
	"https://spacetime.dev/papers/minimal-surface-theory.pdf",
}

const sampleHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#000000" />
<title>home :: spacetime.dev</title>
<link rel="stylesheet" href="/assets/styles/main.min.css">
<link rel="alternate" type="application/atom+xml" href="https://spacetime.dev/feed.xml" />
</head>
<body>
<div class="menu">
[<a href="/">home</a>:
:<a href="/posts">posts</a>:
:<a href="/feed.xml">subscribe</a>]
</div>
<div class="main">
<h1 id="awn-umar">awn umar</h1>
<div class="nav">
<a href="/cdn-cgi/l/email-protection#1372647d536063727076677a7e763d777665" rel="me">email</a> : : <a href="https://spacetime.dev/public-key.pgp" rel="me">pgp</a> : : <a href="https://github.com/awnumar" rel="me">github</a> : : <a href="https://read.cv/awn" rel="me">cv</a>
</div>
<p>Programmer mostly <a href="https://github.com/awnumar">working on</a> security stuff. Mathematics and computer science at the <a href="https://en.wikipedia.org/wiki/University_of_Bristol">University of Bristol</a>. Here I <a href="/posts">post</a> about things I’m working on or thinking about.</p>
<h2 id="recent-posts">recent posts</h2>
<div class="nav">
<a href="/posts">all posts</a> : : <a href="/feed.xml">subscribe</a>
</div>
<ul>
<li>2020-12-16 : : <a href="/rosen-censorship-resistant-proxy-tunnel">rosen: censorship-resistant proxy tunnel</a></li>
<li>2020-02-20 : : <a href="/plausibly-deniable-encryption">plausibly deniable encryption</a></li>
<li>2019-07-20 : : <a href="/memory-retention-attacks">memory retention attacks</a></li>
<li>2019-07-18 : : <a href="/encrypting-secrets-in-memory">encrypting secrets in memory</a></li>
</ul>
<h2 id="projects">projects</h2>
<ul>
<li><a href="https://github.com/awnumar/memguard">memguard</a> : : secure software enclave for storage of sensitive information in memory</li>
<li><a href="https://github.com/awnumar/rosen">rosen</a> : : modular proxy tunnel that circumvents censorship by encapsulating traffic within a cover protocol</li>
<li><a href="https://github.com/awnumar/gravity">gravity</a> : : experimental user-space data protection utility providing plausibly deniable encryption</li>
</ul>
<h2 id="papers">papers</h2>
<blockquote>
<p>Minimal Surface Theory: The Mathematics of Soap Films [<a href="/papers/minimal-surface-theory.pdf">pdf</a>]</p>
<p>An introduction to minimal surfaces and their relation to soap films, and some ways that they can be created.</p>
</blockquote>
</div>
<script data-cfasync="false" src="/cdn-cgi/scripts/5c5dd728/cloudflare-static/email-decode.min.js"></script><script defer src='https://static.cloudflareinsights.com/beacon.min.js' data-cf-beacon='{"token": "6eaa5a829d7a4c57a0c04c420e1743b0"}'></script>
</body>
</html>`

// note: don't use this on sensitive data, there's a timing side channel attack
func equal(a []*url.URL, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].String() != b[i] {
			return false
		}
	}
	return true
}
