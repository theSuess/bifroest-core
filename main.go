package main

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
	"net/http"
	"net/http/httputil"
	"strings"
)

type handler struct {
	store URLStore
}

type configuration struct {
	BindAddress string
}

// URLStore Abstracts the data storage for easier testing
type URLStore interface {
	Get(string) (string, error)
}

// RedisURLStore implements the URLStore Interface
type RedisURLStore struct {
	client *redis.Client
}

// Get is implemented as specified in the URLStore interface
func (r RedisURLStore) Get(url string) (string, error) {
	return r.client.Get(url).Result()
}

func (h *handler) proxy(r *http.Request) {
	domain := strings.Split(r.Host, ".")[0]

	log.WithFields(log.Fields{
		"domain":     domain,
		"remoteAddr": r.RemoteAddr,
	}).Infof("%s %s", r.Method, r.RequestURI)
	host, _ := h.store.Get("bfr:domains:" + domain)
	r.Host = host
	r.URL.Host = r.Host
	r.URL.Scheme = "https" //TODO: check for correct scheme

}

func main() {
	var config configuration
	if _, err := toml.DecodeFile("bifroest.toml", &config); err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.DebugLevel)

	client := RedisURLStore{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}

	h := handler{store: client}
	prx := &httputil.ReverseProxy{
		Director: h.proxy,
	}
	log.Printf("Server started on `%s`", config.BindAddress)
	log.Fatal(http.ListenAndServe(config.BindAddress, prx))
}
