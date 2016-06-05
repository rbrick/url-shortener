package storage

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

var (
	ErrorFailedToDelete = errors.New("Failed to delete ID")
	ErrorDoesNotExist   = errors.New("Requested ID does not exist")
	ErrorFailedToCreate = errors.New("Failed to create short URL")
	ErrorNotURL         = errors.New("Input is not a URL!")
	ErrorFailedToLookup = errors.New("Failed to lookup")
)

// ShortUrl represents a short url. It contains the time it was created, the id, and the long url
type ShortUrl struct {
	Id        string `json:"id,omitempty"`
	Hash      string `json:"hash,omitempty"`
	LongUrl   string `json:"longUrl,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

func NewShortURL(longUrl string) ShortUrl {
	base := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(longUrl))))
	base64.StdEncoding.Encode(base, []byte(longUrl))
	hash := fmt.Sprintf("%02x", sha1.Sum(base))
	return ShortUrl{hash[0:8], hash, longUrl, time.Now().Unix()}
}

func NewCustomShortURL(id, longUrl string) ShortUrl {
	base := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(longUrl))))
	base64.StdEncoding.Encode(base, []byte(longUrl))
	hash := fmt.Sprintf("%02x", sha1.Sum(base))
	return ShortUrl{id, hash, longUrl, time.Now().Unix()}
}

// ShortenService abstracts interactions between the database and the API
type ShortenService interface {
	// Lookup performs a query on the database for a short url
	Lookup(id string) (*ShortUrl, error)
	// ReverseLookup lookups up a short URL in the database from the longUrl
	ReverseLookup(longUrl string) (*ShortUrl, error)
	// Exists checks if the id exists
	Exists(id string) bool
	// Insert creates a short url
	Insert(*ShortUrl) error
	// Delete deletes a short url
	Delete(id string) error
}
