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

const URL_PATTERN = `^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`

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

	// Handle API routes
	router.HandleFunc("/api/lookup", ApiLookupHandler)
	router.HandleFunc("/api/create", ApiCreateHandler)

	// Handle static files
	router.PathPrefix("/").Handler(fileServer{"static", true})

	// Handle 404 errors
	router.NotFoundHandler = notFoundHandler{}

	server := server{router}

	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), server)
}

type server struct {
	delegate http.Handler
}

func (s server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "-", req.URL.Path)
	s.delegate.ServeHTTP(res, req)
}
