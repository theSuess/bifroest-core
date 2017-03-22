package main

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

type dummyStore struct{}

func (d dummyStore) Get(url string) (string, error) {
	switch url {
	case "foo":
		return "http://bar.com", nil
	case "bar":
		return "https://foo.com", nil
	default:
		return "", errors.New("url not found")
	}
}

func TestProxy(t *testing.T) {
	h := handler{
		store:     dummyStore{},
		errorPage: "https://example.com",
	}
	r := &http.Request{}
	r.Host = "foo"
	r.URL = &url.URL{}
	r.URL.Host = "foo"
	h.proxy(r)
	if r.Host != "http://bar.com" && r.URL.Host != r.Host {
		t.Errorf("Proxy rewrite did not work! Expected host to be %s but was %s", "http://bar.com", r.Host)
	}
	r.Host = "bar"
	r.URL.Host = "bar"
	h.proxy(r)
	if r.Host != "https://foo.com" && r.URL.Host != r.Host {
		t.Errorf("Proxy rewrite did not work! Expected host to be %s but was %s", "https://foo.com", r.Host)
	}
	r.Host = "baz"
	r.URL.Host = "baz"
	h.proxy(r)
	if r.Host != "https://example.com" && r.URL.Host != r.Host {
		t.Errorf("Proxy rewrite did not work! Expected host to be %s (errorPage) but was %s", "https://example.com", r.Host)
	}
}
