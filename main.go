package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type handler struct {
	client *redis.Client
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Host + r.RequestURI)
	domain := strings.Split(r.Host, ".")[0]
	host, err := h.client.Get("bfr:domains:" + domain).Result()
	if err != nil {
		log.Warnf("Invalid Subdomain requested: %s", domain)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Subdomain not found"))
		return
	}

	uri := host + r.RequestURI

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
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	h := handler{client: client}
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      &h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("Server started")
	log.Fatal(srv.ListenAndServe())
}
