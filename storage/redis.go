package storage

import (
	"github.com/garyburd/redigo/redis"
)

const REDIS_SHORTURLS = "shorturls:"
const REDIS_HASHES = "shorturls_hashes:"

type RedisShortenService struct {
	urlCache  map[string]ShortUrl
	hashCache map[string]string
	conn      redis.Conn
}

// performs a lookup on the database for a short url
func (r *RedisShortenService) Lookup(id string) (*ShortUrl, error) {
	if v, ok := r.urlCache[id]; ok {
		return &v, nil
	}

	if r.Exists(id) {
		// Send HGETALL Redis command, and get the string values as interfaces which we can use
		val, err := redis.Values(r.conn.Do("HGETALL", REDIS_SHORTURLS+id))
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

		r.urlCache[id] = res
		r.hashCache[res.Hash] = id
		return &res, nil
	}
	return nil, ErrorDoesNotExist
}

func (r *RedisShortenService) Exists(id string) bool {
	return r.exists("url", id)
}

func (r *RedisShortenService) Delete(id string) error {
	return ErrorFailedToDelete
}

func (r *RedisShortenService) Insert(url *ShortUrl) error {
	if r.Exists(url.Id) {
		return nil
	}
	_, err := r.conn.Do("HMSET", REDIS_SHORTURLS+url.Id, "Id", url.Id, "LongUrl", url.LongUrl, "Hash", url.Hash, "Timestamp", url.Timestamp)
	if err != nil {
		return err
	}

	_, err = r.conn.Do("SET", REDIS_HASHES+url.Hash, url.Id)
	if err != nil {
		return err
	}

	r.urlCache[url.Id] = *url
	r.hashCache[url.Hash] = url.Id
	return nil
}

func (r *RedisShortenService) ReverseLookup(longUrl string) (*ShortUrl, error) {
	query := NewShortURL(longUrl)
	if r.exists("hash", query.Hash) {
		if v, ok := r.hashCache[query.Hash]; ok {
			return r.Lookup(v)
		}

		v, err := redis.String(r.conn.Do("GET", REDIS_HASHES+query.Hash))
		if err != nil {
			return nil, err
		}
		r.hashCache[query.Hash] = v
		return r.Lookup(v)
	}
	return nil, ErrorFailedToLookup
}

func (r *RedisShortenService) exists(lookupType, id string) bool {
	var key string

	switch lookupType {
	case "url":
		{
			if _, ok := r.urlCache[id]; ok {
				return ok
			}
			key = REDIS_SHORTURLS
		}
	case "hash":
		if _, ok := r.hashCache[id]; ok {
			return ok
		}
		key = REDIS_HASHES
	default:
		return false
	}

	res, err := redis.Bool(r.conn.Do("EXISTS", key+id))
	if err != nil {
		return false
	}
	return res
}

func NewRedisShortenService(conn redis.Conn) *RedisShortenService {
	return &RedisShortenService{map[string]ShortUrl{}, map[string]string{}, conn}
}
