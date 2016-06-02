package storage

import "errors"

var (
	ErrorFailedToDelete = errors.New("Failed to delete ID")
	ErrorDoesNotExist   = errors.New("Requested ID does not exist")
	ErrorFailedToCreate = errors.New("Failed to create short URL")
	ErrorFailedToLookup = errors.New("Failed to lookup")
)

// ShortUrl represents a short url. It contains the time it was created, the id, and the long url
type ShortUrl struct {
	Id        string `json:"id,omitempty"`
	LongUrl   string `json:"longUrl,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
}

// ShortenService abstracts interactions between the database and the API
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
