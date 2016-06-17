package storage

import (
	"errors"
)

var (
	ErrorFailedToDelete = errors.New("Failed to delete ID")
	ErrorDoesNotExist   = errors.New("Requested ID does not exist")
	ErrorFailedToCreate = errors.New("Failed to create short URL")
	ErrorNotURL         = errors.New("Input is not a URL!")
	ErrorFailedToLookup = errors.New("Failed to lookup")
)

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

type UserService interface {
	Get(id string) *RegisteredUser
}
