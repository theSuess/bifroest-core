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
		return "bar.com", nil
	case "bar":
		return "foo.com", nil
	default:
		return "", errors.New("url not found")
	}
}

func TestProxy(t *testing.T) {
	h := handler{
		store: dummyStore{},
	}
	r := &http.Request{}
	r.Host = "foo"
	r.URL = &url.URL{}
	r.URL.Host = "foo"
	h.proxy(r)
	if r.Host != "bar.com" && r.URL.Host != r.Host {
		t.Errorf("Proxy rewrite did not work! Expected host to be %s but was %s", "bar.com", r.Host)
	}
	r.Host = "bar"
	r.URL.Host = "bar"
	h.proxy(r)
	if r.Host != "foo.com" && r.URL.Host != r.Host {
		t.Errorf("Proxy rewrite did not work! Expected host to be %s but was %s", "foo.com", r.Host)
	}
}
