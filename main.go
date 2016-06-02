package main

import (
	"./storage"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type Config struct {
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Server struct {
		Port uint16 `yaml:"port"`
	}
}

var (
	config         Config
	shortenService *storage.RedisShortenService
)

func init() {

	// Load the config
	if d, err := ioutil.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	} else {
		err := yaml.Unmarshal(d, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	// At this point we can assume the config has been loaded, soooo connect to the redis
	if r, err := redis.Dial("tcp", net.JoinHostPort(config.Redis.Host, config.Redis.Port)); err != nil {
		log.Fatal(err)
	} else {
		shortenService = storage.NewRedisShortenService(r)
	}
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Handle the redirect functionality
	router.HandleFunc("/{id:[a-zA-Z0-9_]+}", RedirectHandler)

	// Handle static files
	router.PathPrefix("/static/").Handler(fileServer{"static", true})

	// Handle 404 errors
	router.NotFoundHandler = notFoundHandler{}

	// Handle API routes
	router.HandleFunc("/api/lookup", ApiLookupHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router)
}
