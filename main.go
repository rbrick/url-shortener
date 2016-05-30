package main

import (
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
	redisConn      redis.Conn
	config         Config
	shortenService RedisShortenService
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
		redisConn = r
		shortenService = RedisShortenService{map[string]ShortUrl{}}
	}
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Handle the redirect functionality
	router.HandleFunc("/r/{id:[a-zA-Z0-9_]+}", RedirectHandler)
	// Handle the index page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router)
}

// Handles redirecting to a long url
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	log.Println("Requested Redirect")
	id := v["id"]
	log.Println("Requested ID:", id)
	s, err := shortenService.Lookup(id)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, s.LongUrl, http.StatusMovedPermanently)
}
