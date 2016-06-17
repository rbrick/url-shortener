package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rbrick/url-shortener/middleware"
	"github.com/rbrick/url-shortener/storage"
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
		Port            uint16 `yaml:"port"`
		UseHTTPS        bool   `yaml:"use-https"`
		StaticDirectory string `yaml:"static-directory"`
	}
	HTTPS struct {
		KeyFile  string `yaml:"key-file"`
		CertFile string `yaml:"cert-file"`
	}
}

var (
	config         Config
	shortenService *storage.RedisShortenService
	boltDatabase   *bolt.DB
)

const URL_PATTERN = `^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`

func init() {
	log.Println("Loading config.yml...")
	// Load the config
	if d, err := ioutil.ReadFile("config.yml"); err != nil {
		log.Fatal(err)
	} else {
		err := yaml.Unmarshal(d, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Config loaded.")

	log.Println("Setting up service.")
	// At this point we can assume the config has been loaded, soooo connect to the redis
	if r, err := redis.Dial("tcp", net.JoinHostPort(config.Redis.Host, config.Redis.Port)); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to Redis.")
		shortenService = storage.NewRedisShortenService(r)
	}

	log.Println("Done.")
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

	auth := middleware.NewAuthHandler(server)
	auth.AddPath("/me")

	addr := fmt.Sprintf(":%d", config.Server.Port)
	if config.Server.UseHTTPS {
		log.Fatal(http.ListenAndServeTLS(addr, config.HTTPS.CertFile, config.HTTPS.KeyFile, auth))
	} else {
		log.Fatal(http.ListenAndServe(addr, auth))
	}

}

// This wraps all incoming HTTP requests and logs them.
// Very useful.
type server struct {
	delegate http.Handler
}

func (s server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "-", req.URL.Path)
	s.delegate.ServeHTTP(res, req)
	fmt.Println("Status:", res.Header().Get("Status"), res.WriteHeader())
}
