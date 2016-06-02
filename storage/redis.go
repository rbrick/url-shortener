package storage

import (
	"github.com/garyburd/redigo/redis"
)

type RedisShortenService struct {
	cache map[string]ShortUrl
	conn  redis.Conn
}

// performs a lookup on the database for a short url
func (r *RedisShortenService) Lookup(id string) (*ShortUrl, error) {
	if v, ok := r.cache[id]; ok {
		return &v, nil
	}

	if r.Exists(id) {
		// Send HGETALL Redis command, and get the string values as interfaces which we can use
		val, err := redis.Values(r.conn.Do("HGETALL", id))
		if err != nil {
			return nil, err
		}

		var res ShortUrl
		// Set the fields in the ShortUrl struct (Use a pointer so we are not just copying the value)
		err = redis.ScanStruct(val, &res)
		if err != nil {
			return nil, err
		}

		// If the Id does not exist it is either invalid, or doesn't exist
		if res.Id == "" {
			return nil, ErrorDoesNotExist
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
	res, err := redis.Bool(r.conn.Do("EXISTS", id))
	if err != nil {
		return false
	}
	return res
}

func (r *RedisShortenService) Delete(id string) error {
	return ErrorFailedToDelete
}

func (r *RedisShortenService) Insert(id, longUrl string) error {
	return ErrorFailedToCreate
}

func NewRedisShortenService(conn redis.Conn) *RedisShortenService {
	return &RedisShortenService{map[string]ShortUrl{}, conn}
}
