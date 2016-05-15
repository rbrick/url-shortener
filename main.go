package main

import (
	"net/http"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"errors"
	"log"
	"strconv"
)

var redisConn redis.Conn

// A short url, contains the time it was created,
type ShortUrl struct {
	Id        string `json:"id"`
	LongUrl   string `json:"longUrl"`
	Timestamp int64 `json:"timestamp"`
	Clicks    int `json:"clicks"`
}

// performs a lookup on the database for a short url
func lookup(id string) (*ShortUrl, error) {
	val, err := redis.ByteSlices(redisConn.Do("HMGET", id, "Id", "LongUrl", "Timestamp", "Clicks"))
	if err != nil {
		return nil, errors.New("Invalid ID!")
	}
	longUrl := string(val[1])
	timestamp, _ := strconv.Atoi(string(val[2]))
	clicks, _ := strconv.Atoi(string(val[3]))
	return &ShortUrl{
		id,
		longUrl,
		int64(timestamp),
		clicks,
	}, nil
}

func init() {
	if r, err := redis.Dial("tcp", "localhost:6379"); err != nil {
		log.Fatal(err)
	} else {
		redisConn = r
	}
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
	s, err := lookup(v["id"])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, s.LongUrl, http.StatusMovedPermanently)
}