package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct{}

var subdomains = map[string]string{
	"foo": "http://example.com",
	"bar": "http://stackoverflow.com",
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	base := subdomains[strings.Split(r.Host, ".")[0]]

	uri := base + r.RequestURI

	rr, err := http.NewRequest(r.Method, uri, r.Body)
	fatal(err)
	copyHeader(r.Header, &rr.Header)

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	fatal(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(body)
	http.NotFound(w, r)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}

func main() {
	h := Handler{}
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      &h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
