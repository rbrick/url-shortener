package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
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
	redisConn      redis.Conn
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
		redisConn = r
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

	addr := fmt.Sprintf(":%d", config.Server.Port)
	if config.Server.UseHTTPS {
		log.Fatal(http.ListenAndServeTLS(addr, config.HTTPS.CertFile, config.HTTPS.KeyFile, server))
	} else {
		log.Fatal(http.ListenAndServe(addr, server))
	}

}

type server struct {
	delegate http.Handler
}

func (s server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// This path dumps the redis values into JSON.
	// the path does not need to be safe. Errors are ignored.
	if req.URL.Path == "/_redisdump" {
		// a empty slice of short urls
		urls := []storage.ShortUrl{}
		resp, _ := redis.Strings(redisConn.Do("KEYS", "shorturls:*"))

		for _, v := range resp {
			var x storage.ShortUrl
			vals, _ := redis.Values(redisConn.Do("HGETALL", v))
			redis.ScanStruct(vals, &x)
			urls = append(urls, x)
		}

		d, _ := json.Marshal(urls)
		res.Write(d)
		return
	}
	s.delegate.ServeHTTP(res, req)
	log.Println(req.Method, "-", req.URL.Path)
}
