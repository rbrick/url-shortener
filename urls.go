package main

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

var (
	ErrorFailedToDelete = errors.New("Failed to delete ID")
	ErrorDoesNotExist   = errors.New("Requested ID does not exist")
)

// A short url, contains the time it was created,
type ShortUrl struct {
	Id        string `redis:"Id"`
	LongUrl   string `redis:"LongUrl"`
	Timestamp int    `redis:"Timestamp"`
}

type ShortenService interface {
	// Lookup performs a query on the database for a short url
	Lookup(id string) (*ShortUrl, error)
	// Exists checks if the id exists
	Exists(id string) bool
	// Insert creates a short url
	Insert(id, longUrl string) error
	// Delete deletes a short url
	Delete(id string) error
}

type RedisShortenService struct {
	cache map[string]ShortUrl
}

// performs a lookup on the database for a short url
func (r *RedisShortenService) Lookup(id string) (*ShortUrl, error) {
	if v, ok := r.cache[id]; ok {
		return &v, nil
	}

	if r.Exists(id) {
		// Send HGETALL Redis command, and get the string values as interfaces which we can use
		val, err := redis.Values(redisConn.Do("HGETALL", id))
		if err != nil {
			return nil, err
		}

		var res ShortUrl
		// Set the fields in the ShortUrl struct (Use a pointer so we are not just copying the value)
		err = redis.ScanStruct(val, &res)
		if err != nil {
			return nil, err
		}
		r.cache[id] = res
		return &res, nil
	}
	return nil, ErrorDoesNotExist
}

func (r *RedisShortenService) Exists(id string) bool {
	if _, ok := r.cache[id]; ok {
		return ok
	}
	res, err := redis.Bool(redisConn.Do("EXISTS", id))
	if err != nil {
		return false
	}
	return res
}

func (r *RedisShortenService) Delete(id string) error {
	return ErrorFailedToDelete
}
