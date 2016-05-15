package main

import (
	"time"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var redisConn redis.Conn

// A short url, contains the time it was created,
type ShortUrl struct {
	Id string `json:"id"`
	LongUrl string `json:"longUrl"`
	Timestamp int64 `json:"timestamp"`
	Clicks int `json:"clicks"`
}

// Creates a new short url
func create(id, longUrl string) ShortUrl {
	return ShortUrl{id, longUrl, time.Now().Unix(), 0}
}

func lookup(id string) {
}


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{id:[\\w]+}", redirect)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":8080", router)
}

// Handles redirecting to a long url
func redirect(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	w.Write([]byte(v["id"]))
}